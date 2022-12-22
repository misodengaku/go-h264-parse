// generated by deep-copy -o sei_deepcopy.go --type SEI .; DO NOT EDIT.

package h264parse

// DeepCopy generates a deep copy of SEI
func (o SEI) DeepCopy() SEI {
	var cp SEI = o
	if o.PayloadBytes != nil {
		cp.PayloadBytes = make([]byte, len(o.PayloadBytes))
		copy(cp.PayloadBytes, o.PayloadBytes)
	}
	return cp
}
