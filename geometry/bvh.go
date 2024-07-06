package geometry

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

func BuildBVH(prims []Primitive) *BVHNode {
	n := len(prims)

	var left Primitive
	var right Primitive
	var bbox = EmptyAABB()

	if n == 1 {
		// leaf node case
		left = prims[0]
		right = nil
		bbox = left.Bounds()
	} else if n == 2 {
		left = prims[0]
		right = prims[1]
		bbox = left.Bounds().Union(right.Bounds())
	} else {
		// interior node case

		// 1. compute a compound bounds of all primitives
		for _, o := range prims {
			bbox = bbox.Union(o.Bounds())
		}

		// 2. calculate bounds for the centroids and pick the longest axis
		var centroidBounds = EmptyAABB()
		for _, p := range prims {
			centroidBounds = centroidBounds.UnionPoint3(centroid(p))
		}
		axis := centroidBounds.LongestAxis()

		// 3. find a midpoint of the centroids
		c1, c2 := centroidBounds.GetCoordinatesByAxis(axis)
		mPoint := (c1 + c2) / 2

		// 4. divide set of primitives into two equal parts such that coordinate of a centroid < pMid goes to the first half, and other goes to the other
		mid := partition(prims, func(p Primitive) bool {
			return centroid(p).GetCoordinateByAxis(axis) < mPoint
		})

		if mid == n {
			mid = n / 2
		}

		left = BuildBVH(prims[0:mid])
		right = BuildBVH(prims[mid:n])
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

func (n *BVHNode) Bounds() Bounds3 {
	return n.Box
}

func centroid(p Primitive) Point3 {
	primBounds := p.Bounds()
	return primBounds.Pmin.Scale(.5).Add(primBounds.Pmax.Scale(.5))
}
