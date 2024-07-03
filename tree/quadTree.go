package tree

import (
	"particle"
	"sync"
)

const N_CHILD int = 4

type QuadTree struct {
	particle                 *particle.Particle // Store the particle
	count                    int                // Number of particles in the subtree - asssuming unit mass, can be used as weight
	lowX, lowY, highX, highY float64            // Coordinates that the subtree contains
	parent                   *QuadTree          // Parent of the node
	child                    [N_CHILD]*QuadTree // Children of the node
	lock                     *sync.Mutex
}

// Create a new empty node
func CreateNode(par *QuadTree, lowX, lowY, highX, highY float64) *QuadTree {
	var node *QuadTree = &QuadTree{}
	node.particle = nil
	node.parent = par
	for i := 0; i < N_CHILD; i++ {
		node.child[i] = nil
	}
	node.lowX = lowX
	node.lowY = lowY
	node.highX = highX
	node.highY = highY
	node.count = 0
	node.lock = &sync.Mutex{}
	return node
}

// Add particle to tree
func AddParticleToTree(p *particle.Particle, node *QuadTree) {
	node.count = node.count + 1
	if node.count == 1 { // Node does not have a particle
		node.particle = p // Add the particle to the node
	} else {
		// If node contains a particle, remove and add the particle to a child subtree
		if node.particle != nil {
			DecideChildSubtree(node.particle, node)
			node.particle = nil
		}
		DecideChildSubtree(p, node)
	}
}

// Add particle to the subtree of the corresponding child
func DecideChildSubtree(p *particle.Particle, node *QuadTree) {
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
	// Add particle to subtree of corresponding node
	if p.X <= midX {
		if p.Y <= midY {
			AddParticleToTree(p, node.child[0])
		} else {
			AddParticleToTree(p, node.child[2])
		}
	} else {
		if p.Y <= midY {
			AddParticleToTree(p, node.child[1])
		} else {
			AddParticleToTree(p, node.child[3])
		}
	}
}

// Calculate force on a particle using the quad tree
func CalcTreeForce(p *particle.Particle, node *QuadTree, theta float64, dt float64) {
	if node.count == 1 {
		particle.CalcForce(p, node.particle, dt)
	} else {
		var r float64 = particle.CalcDistance(p, node.particle)
		var D float64 = node.highX - node.lowX
		if D/r < theta {
			particle.CalcForce(p, node.particle, dt)
		} else {
			for i := 0; i < N_CHILD; i++ {
				if node.child[i] != nil && node.child[i].count > 0 {
					CalcTreeForce(p, node.child[i], theta, dt)
				}
			}
		}
	}
}

// Compute center of mass of a particle
func ComputeCenterOfMass(node *QuadTree) {
	if node == nil || node.count <= 1 {
		return
	}
	for i := 0; i < N_CHILD; i++ {
		ComputeCenterOfMass(node.child[i])
	}
	var cmX float64 = 0.0
	var cmY float64 = 0.0
	var count int = 0
	for i := 0; i < N_CHILD; i++ {
		if node.child[i].count > 0 {
			cmX = cmX + float64(node.child[i].count)*node.child[i].particle.X
			cmY = cmY + float64(node.child[i].count)*node.child[i].particle.Y
			count = count + node.child[i].count
		}
	}
	node.particle = &particle.Particle{}
	node.particle.X = cmX / float64(count)
	node.particle.Y = cmY / float64(count)
	node.particle.Vx = 0
	node.particle.Vy = 0
}
