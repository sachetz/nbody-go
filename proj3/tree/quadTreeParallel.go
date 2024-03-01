package tree

import (
	"math"
	"proj3/barrier"
	"proj3/particle"
	"proj3/utils"
	"sync"
)

func AddParticlesParallel(p []*particle.Particle, root *QuadTree, nParticles int, numThreads int) {
	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		lowerBound := i * nParticles / numThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/numThreads), float64(nParticles)))
		f := func(tid int) {
			defer wg.Done()
			for j := lowerBound; j < upperBound; j++ {
				AddParticleToTreeParallel(p[j], root, false)
			}
		}
		go f(i)
	}
	wg.Wait()
}

func chooseNode(p *particle.Particle, node *QuadTree, midX, midY float64, update_parent_after bool) {
	if p.X <= midX {
		if p.Y <= midY {
			AddParticleToTreeParallel(p, node.child[0], update_parent_after)
		} else {
			AddParticleToTreeParallel(p, node.child[2], update_parent_after)
		}
	} else {
		if p.Y <= midY {
			AddParticleToTreeParallel(p, node.child[1], update_parent_after)
		} else {
			AddParticleToTreeParallel(p, node.child[3], update_parent_after)
		}
	}
}

// Add particle to tree
func AddParticleToTreeParallel(p *particle.Particle, node *QuadTree, update_parent bool) {
	node.lock.Lock()
	if node.parent != nil && !update_parent {
		node.parent.lock.Unlock()
	}
	node.count = node.count + 1
	if node.count == 1 { // Node does not have a particle
		node.particle = p // Add the particle to the node
		node.lock.Unlock()
	} else {
		// If node contains a particle, remove and add the particle to a child subtree
		var midX, midY float64
		midX = (node.lowX + node.highX) / 2
		midY = (node.lowY + node.highY) / 2
		if node.child[0] == nil { // Node does not have children
			// Create children for node
			node.child[0] = CreateNode(node, node.lowX, node.lowY, midX, midY)   // Lower left
			node.child[1] = CreateNode(node, midX, node.lowY, node.highX, midY)  // Upper left
			node.child[2] = CreateNode(node, node.lowX, midY, midX, node.highY)  // Lower Right
			node.child[3] = CreateNode(node, midX, midY, node.highX, node.highY) // Upper Right
		}
		if node.particle != nil {
			chooseNode(node.particle, node, midX, midY, true)
			node.particle = nil
		}
		chooseNode(p, node, midX, midY, false)
	}
}

func CalcTreeForceAndUpdatePosParallel(p []*particle.Particle, root *QuadTree, numThreads int, nParticles int) {
	// Calculate force on each particle
	var wg sync.WaitGroup
	b := barrier.NewBarrier(numThreads)
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		lowerBound := i * nParticles / numThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/numThreads), float64(nParticles)))
		f := func(tid int, bar *barrier.Barrier) {
			defer wg.Done()
			for j := lowerBound; j < upperBound; j++ {
				CalcTreeForce(p[j], root, utils.Theta, utils.Dt)
			}
			bar.Wait()
			particle.UpdatePosSequential(p, lowerBound, upperBound)
		}
		go f(i, b)
	}
	wg.Wait()
}
