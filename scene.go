package fauxgl

import (
	"log"
	"sync"
)

// Scene struct to store all data for a scene
type Scene struct {
	Context *Context
	Objects []*Object
}

// NewScene returns new scene
func NewScene(context *Context) *Scene {
	return &Scene{context, nil}
}

// AddObject adds object to scene
func (s *Scene) AddObject(o *Object) {
	s.Objects = append(s.Objects, o)
}

// FitObjectsToScene fit the objects into a 0.5 unit bounding box
func (s *Scene) FitObjectsToScene(eye, center, up Vector, fovy, aspect, near, far float64) (matrix Matrix) {
	matrix = LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := NewPhongShader(matrix, Vector{}, eye)

	allMesh := NewEmptyMesh()
	var boxes []Box
	for _, o := range s.Objects {
		if o.Mesh == nil {
			continue
		}
		allMesh.Add(o.Mesh)
		bb := o.Mesh.BoundingBox()
		boxes = append(boxes, bb)
	}
	box := BoxForBoxes(boxes)
	b := NewCubeForBox(box)
	b.BiUnitCube()
	allMesh.FitInside(b.BoundingBox(), V(0.5, 0.5, 0.5))

	indexed := 0
	var addedFOV float64
	for _, o := range s.Objects {
		if o.Mesh == nil {
			continue
		}
		num := len(o.Mesh.Triangles)
		tris := allMesh.Triangles[indexed : num+indexed]
		allInside := false
		for !allInside && len(tris) > 0 {
			for _, t := range tris {
				v1 := shader.Vertex(t.V1)
				v2 := shader.Vertex(t.V2)
				v3 := shader.Vertex(t.V3)

				if v1.Outside() || v2.Outside() || v3.Outside() {
					addedFOV += 5
					matrix = LookAt(eye, center, up).Perspective(fovy+addedFOV, aspect, near, far)
					shader.Matrix = matrix
					allInside = false
				} else {
					allInside = true
				}
			}
		}

		o.Mesh = NewTriangleMesh(tris)
		indexed += num
	}

	return
}

// Draw draws the scene
func (s *Scene) Draw() {
	var wg sync.WaitGroup
	wg.Add(len(s.Objects))
	for _, o := range s.Objects {
		if o.Mesh == nil {
			wg.Done()
			log.Printf("Objected attempted to render with nil mesh")
			continue
		}
		go s.Context.DrawObject(o, &wg)
	}
	wg.Wait()
}
