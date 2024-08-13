package requests

type ScanType string

const (
	SynScan     ScanType = "s"
	ConnectScan ScanType = "c"
)

const (
	// 使用するスレッド数を制御する。
	Threads int = 10
	// ターゲットホストからrespondが返ってくる待機時間
	Timeout int = 30
	// 実行最大時間（分）
	MaxEnumerationTime int = 10
	// クローリングの標準深度
	MaxDepth = 3
	// Delayは、各クロールリクエスト間の遅延時間（秒）
	Delay = 0
	// RateLimitは、1秒間に送信するリクエストの最大数
	RateLimit = 150
)

type FieldScopeType string

const (
	RDN       FieldScopeType = "rdn"
	FQDN      FieldScopeType = "fqdn"
	SUBDOMAIN FieldScopeType = "subdomain"
)

type StrategyType string

const (
	DepthFirst   StrategyType = "depth-first"
	BreadthFirst StrategyType = "breadth-first"
)

// このリクエスト構造体は全ての攻撃を一つのRESTエンドポイントから
// 実行する全ての攻撃のリクエストの形である。
// 個別の攻撃は順次対応する予定。
// 使用されているパラメーターは各構造体に記述している、paramterの詳細を確認してください。
type HackDoScanRequest struct {
	Name string `json:"name"`
	// target parametr formats
	// general domain, IPV4 addresses, etc. are valid; some IPv6 addresses are not supported, so deal with errors if they occur.
	// example:
	// "api.caterpie.runetale.com", "150.155.1.1.1", "runetale.com"
	// *targetとなるホストのURL、全ての攻撃の対象
	TargetHost []string `json:"target"`
	// httpxで使用されるパラメーター
	HTTPMethods string `json:"methods"`
	// portscanはsyn scanとconnect scanに大きく分けられる。
	// syn scanはステルス性が高く、TCPのコネクションを完全に確立しない。ACKを返さない。
	// connect scanはステルス性が低いが、ACKパケットを返すことでTCPのコネクションを確立する。
	// defaultは"s"
	PortScanType ScanType `json:"scan_type"`
	// スキャンするPort
	// フォーマットは"1-500"のように範囲をしてもらう
	// 単一のポートは"22"などで良い、ポートスキャンで使用
	TargetPorts string `json:"ports"`

	// Queriesはcencysで使用する検索クエリ
	Queries []string `json:"queries"`

	// Web Crowling Parameter
	//
	// MaxDepthはクローリング進度を決定する
	// 深度0: クローラーが最初に訪れるページ(スタートページ)
	// 深度1: スタートページから直接リンクされているページ
	// 深度2: 深度1のページからさらにリンクされているページです。
	// 深度3以降: 深度2のページからリンクされているページが深度3、さらにそのリンクからのページが深度4と続く
	// 本来このパラメーターはサービスの特性やvisionアルゴリズムによって、適切に設定されるべきである。
	// 現状は手動指定
	MaxDepth int `json:"max_depth"`

	// Fieldscopeは対象のドメインスコープの範囲を決める
	// "rdn" 登録されているドメイン、example.comのようなドメインが対象
	// "fqdn" 完全修飾ドメイン名を対象とするスコープ。ホスト名やサブドメインを含む完全なドメイン名
	// "subdomain" sub.example.comのようなサブドメインが対象
	// デフォルトは"rdn"
	FieldScope FieldScopeType `json:"field_scope"`

	// タイムアウトはリクエストを待つ時間（秒）
	Timeout int `json:"timeout"`

	// 同時にクロールするゴルーチンの数
	// Concurrencyのようなパラメーターにも同意義で使用される
	Threads int `json:"thread"`

	// 並列性は、URLを処理するゴルーチンの数
	Parallelism int `json:"parallelism"`

	// Delayは、各クロールリクエスト間の遅延時間（秒）
	Delay int `json:"delay"`

	// RateLimitは、1秒間に送信するリクエストの最大数
	RateLimit int `json:"rate_limit"`

	// リトライ回数
	Retry int `json:"retry"`

	// クローリングの戦略を指定できる。
	// "depth-first"と"breadth-first"の2つの戦略を選択できる。
	// Depth-Firstの特徴
	// 	- 特定のコンテンツを迅速に探索: 特定の階層にあるコンテンツに迅速にアクセスしたい場合に有効
	//  - メモリ効率が高い: 常に現在のパスのみを保持するため、メモリ使用量が少ない
	//  - リソースの偏り: 特定のパスにリソースを集中させるため、偏りが生じる可能性がある
	// Breadth-Firstの特徴
	// 	- 全体像の把握: サイト全体の構造を把握しやすく、特定の階層内の全リンクを均等に探索可能
	//  - メモリ使用量が多い: 各階層の全リンクを保持するため、メモリ使用量が多くなる
	//  - 初期探索が遅い: 広範囲に探索するため、特定の深い階層に到達するまでに時間がかかる
	// デフォルトは"depth-frist"
	Strategy StrategyType `json:"strategy"`

	// nucleiを使用したテンプレート攻撃の種類を選択できる。
	// デフォルトではこの量のtempateを使用して攻撃をする
	// TemplteTags: []string{
	// 	"cve", "cnvd", "nvd", "owasp", "fuzz", "misc", "oast", "default",
	// 	"network", "ssl", "technologies", "takeovers", "dns", "files",
	// 	"tokens", "vulnerabilities", "web", "wordpress", "waf", "miscellaneous",
	// },
	// https://github.com/projectdiscovery/nuclei-templates ここで定義されている、idを使用できる。
	TemplateTags []string `json:"template"`
}
