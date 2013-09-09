
package proto

type Serializer interface {
	func Serialize() ([]byte, error)
	func Deserialize(packet []byte) error
}