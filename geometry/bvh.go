package geometry

import (
	"github.com/danradchuk/raytracer/shading"
)

type BVHNode struct {
	Left  Primitive
	Box   Bounds3
	Right Primitive
}

// To construct a BVH we need
// 1. Compute a AABB for every triangle of the mesh
// 2. Compute a centroid of the AABB
// 3. Compute a midpoint of the centroids
// 4. Split a slice of Primitives by the midpoint
// 5. Recursively build a BVH

func BuildBVH(objects []Primitive) *BVHNode {
	n := len(objects)

	var left Primitive
	var right Primitive
	var bbox = EmptyAABB()

	if n == 1 {
		// leaf node case
		left = objects[0]
		right = nil
		bbox = left.BBox()
	} else if n == 2 {
		left = objects[0]
		right = objects[1]
		bbox = left.BBox().Union(right.BBox())
	} else {
		// interior node case

		// 1. compute a compound bounds of all primitives in objects
		for _, o := range objects {
			bbox = bbox.Union(o.BBox())
		}

		// 2. calculate bounds for the centroids and pick the longest axis
		var centroidBounds = EmptyAABB()
		for _, o := range objects {
			centroidBounds = centroidBounds.UnionPoint3(o.Centroid())
		}
		axis := centroidBounds.LongestAxis()

		// 3. find a midpoint of the centroids
		c1, c2 := centroidBounds.GetCoordinatesByAxis(axis)
		mPoint := (c1 + c2) / 2

		// 4. divide set of primitives into two equal parts such that coordinate of a centroid < pMid goes to the first half, and other goes to the other
		mid := partition(objects, func(p Primitive) bool {
			return p.Centroid().GetCoordinateByAxis(axis) < mPoint
		})

		if mid == n {
			mid = n / 2
		}

		left = BuildBVH(objects[0:mid])
		right = BuildBVH(objects[mid:n])
	}

	return &BVHNode{
		Left:  left,
		Box:   bbox,
		Right: right,
	}
}

func partition(slice []Primitive, predicate func(Primitive) bool) int {
	i := 0
	j := len(slice) - 1
	for i < j {
		for i < len(slice) && predicate(slice[i]) {
			i++
		}
		for j >= 0 && !predicate(slice[j]) {
			j--
		}
		if i < j {
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
	return i
}

func (n *BVHNode) GetMaterial() shading.Material {
	return shading.Material{}
}

func (n *BVHNode) Intersect(r Ray) *HitRecord {
	var hit *HitRecord = nil
	if n.Box.Intersect(r) {
		var leftHit *HitRecord
		var rightHit *HitRecord

		if n.Left != nil {
			leftHit = n.Left.Intersect(r)
		}
		if n.Right != nil {
			rightHit = n.Right.Intersect(r)
		}

		if leftHit != nil && rightHit != nil {
			if leftHit.T < rightHit.T {
				hit = leftHit
			} else {
				hit = rightHit
			}
		} else if leftHit != nil {
			hit = leftHit
		} else if rightHit != nil {
			hit = rightHit
		}
	}

	return hit
}

func (n *BVHNode) NormalAt(_ Vec3) Vec3 {
	return Vec3{}
}

func (n *BVHNode) BBox() Bounds3 {
	return n.Box
}

func (n *BVHNode) Centroid() Point3 {
	return Point3{}
}
