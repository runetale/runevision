// Engineは攻撃を行ったり、love machineやvisionのアルゴリズムにデータを共有する。
// visionのコアとなる構造体です。
package engine

import "github.com/runetale/runevision/vsd"

// todo:
// apiからパラメーターを受け取るコールバック関数
// apiからvisionを実行するコールバック関数
// この２つを用意する
type Engine struct {
	// vision system
	vsd *vsd.VisionSystem
}
