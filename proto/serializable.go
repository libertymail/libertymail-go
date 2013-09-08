
package proto

type Serializable interface {
	func Serialize() ([]byte, error)
	func Deserialize(packet []byte) error
}