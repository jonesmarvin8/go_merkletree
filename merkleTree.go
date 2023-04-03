package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"math"
)

type MerkleTree struct {
	root       *MerkleNode
	layerCount float64
}

type MerkleNode struct {
	data       string
	dataHash   []byte
	leftChild  *MerkleNode
	rightChild *MerkleNode
	parentNode *MerkleNode
}

type pathMerkle struct {
}

func main() {
	//Test Data
	//path := list.New()

	var inputDATA [8]string
	inputDATA[0] = "Hi"
	inputDATA[1] = "Hello"
	inputDATA[2] = "Hola"
	inputDATA[3] = "Hey"
	inputDATA[4] = "TEST1"
	inputDATA[5] = "TEST2"
	inputDATA[6] = "TEST3"
	inputDATA[7] = "TEST4"

	var num int
	num = 8

	levelCount := math.Log2(float64(num))

	hashFUN := sha256.New()
	tree := &MerkleTree{root: &MerkleNode{data: "", dataHash: hashFUN.Sum(nil), parentNode: nil, leftChild: nil, rightChild: nil}, layerCount: levelCount}

	initMerkleTree(tree.layerCount, tree.root, sha256.New())

	for i := 0; i < num; i++ {
		fillMerkleTree(int(levelCount-1), i, inputDATA[i], tree.root, sha256.New())
	}

	var pathTest [][]byte
	var queryData string
	pathTest, queryData = getMerklePath(tree, 7)

	if checkMerklePath(7, queryData, tree.root.dataHash, pathTest, sha256.New()) {
		fmt.Println("Success")
	}
}

// getMerkleRoot returns the root of a given Merkle Tree.
func (tree *MerkleTree) getMerkleRoot() []byte {
	return tree.root.dataHash
}

// Initializes a blank binary tree of the desired size for levelCount.
// Recursively used to construct subtree.
func initMerkleTree(levelCount float64, currNode *MerkleNode, hashFun hash.Hash) {
	if levelCount > 0 {
		currNode.leftChild = &MerkleNode{data: "", dataHash: hashFun.Sum(nil), leftChild: nil, rightChild: nil, parentNode: nil}
		initMerkleTree(levelCount-1, currNode.leftChild, hashFun)

		currNode.rightChild = &MerkleNode{data: "", dataHash: hashFun.Sum(nil), leftChild: nil, rightChild: nil, parentNode: nil}
		initMerkleTree(levelCount-1, currNode.rightChild, hashFun)
	}
}

// Fill the Merkle tree for leaf location.
// Recursive function to update the hashes after the leaf is reached.
// Assumes that initMerkleTree has been ran first.
func fillMerkleTree(exp int, location int, data string, currNode *MerkleNode, hashFun hash.Hash) {
	if currNode.leftChild == nil {
		currNode.data = data
		hashFun.Write([]byte(data))
		currNode.dataHash = hashFun.Sum(nil)
	} else {
		if location/int(math.Pow(2, float64(exp))) == 0 {
			fillMerkleTree(exp-1, location, data, currNode.leftChild, hashFun)
		} else {
			fillMerkleTree(exp-1, location-int(math.Pow(2, float64(exp))), data, currNode.rightChild, hashFun)
		}

		hashFun.Reset()
		hashFun.Write(currNode.leftChild.dataHash)
		hashFun.Write(currNode.rightChild.dataHash)
		currNode.dataHash = hashFun.Sum(nil)
	}
}

// Given MerkleTree tree and queryLocation, outputs
// the appropriate path to (queryLocation)th leaf and
// the unhashed data associated to the leaf.
func getMerklePath(tree *MerkleTree, queryLocation float64) ([][]byte, string) {
	var path [][]byte
	var currNode *MerkleNode

	path = nil
	currNode = tree.root

	//Out of range query.
	if queryLocation >= math.Pow(2, tree.layerCount) || queryLocation < 0 {
		return path, ""
	}

	for i := 0; i < int(tree.layerCount); i++ {
		if queryLocation-math.Pow(2, tree.layerCount-1-float64(i)) >= 0 {
			path = append(path, currNode.leftChild.dataHash)

			currNode = currNode.rightChild
			queryLocation -= math.Pow(2, tree.layerCount-1-float64(i))
		} else {
			path = append(path, currNode.rightChild.dataHash)

			currNode = currNode.leftChild
		}
	}

	//Reverse order of path for ease of verification.
	var tempByte []byte

	for i := 0; i < 1; i++ {
		tempByte = path[0]
		path[i] = path[len(path)-1-i]
		path[len(path)-1-i] = tempByte
	}

	return path, currNode.data
}

// Given queryData for the queryLocation and the corresponding
// claimed Merkle path. This function computes the appropriate
// hash at each level of the path. At the end, compares the
// resulted computed hash with rootHash.
func checkMerklePath(queryLocation int, queryData string, rootHash []byte, path [][]byte, hashFun hash.Hash) bool {
	var currHash []byte
	hashFun.Write([]byte(queryData))
	currHash = hashFun.Sum(nil)

	for i := 0; i < len(path); i++ {
		hashFun.Reset()

		if queryLocation%2 == 0 {
			hashFun.Write(currHash)
			hashFun.Write(path[i])
			currHash = hashFun.Sum(nil)
		} else {
			hashFun.Write(path[i])
			hashFun.Write(currHash)
			currHash = hashFun.Sum(nil)
		}
		queryLocation /= 2
	}

	if bytes.Compare(rootHash, currHash) == 0 {
		return true
	} else {
		return false
	}
}
