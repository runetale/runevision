package driver

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func processNewFile(fileName string) {
	// 新しいファイルに対して行いたい処理をここに記述
	fmt.Printf("新しいファイルを発見しました: %s\n", fileName)
	// 例: ファイルの中身を読み込む
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("ファイルの読み込みに失敗しました: %v\n", err)
		return
	}
	fmt.Printf("ファイルの内容:\n%s\n", string(data))
}

func watchDirectory(directory string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					// 新しいファイルが作成された場合の処理
					processNewFile(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("エラー:", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}

	// 無限に待機し、ディレクトリの変更を監視し続ける
	<-done
}
