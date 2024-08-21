// vision systemで使う型を定義
// あとから移動させるかもしれない
package types

type SequenceID string

type AttackStatus string

const (
	STARTING  AttackStatus = "STARTING"
	SCANNING  AttackStatus = "SCANNING"
	COMPLETED AttackStatus = "COMPLETED"
)
