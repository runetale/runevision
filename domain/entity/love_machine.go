// love machine entityはmlのために使用する
package entity

// AttackTypeはlove-machine mlが任意のドメインの
// tech stack, port, classなどを学習し
// 選択する単一の攻撃タイプである
type AttackType string

// RT_のプレフィックスが入っている場合は、runetaleにbusiness impactが見込める攻撃手法である
// このサービスは如何にターゲット情報システム部門にサービスが脆弱かを可視化してあげられる。
// runetaleを売り込めるチャンスを秘めている。
const (
	// ** java scriptで外部から実行できる脅威
	// 外部から検出できる mysql audit log
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/javascript/audit/mysql
	JS_EXPLOIT_AUDIT_MYSQL AttackType = "JS_EXPLOIT_AUDIT_MYSQL"

	// backdoor of 'proftpd-1.3.3c'
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/javascript/backdoor
	JS_EXPLOIT_PROFTPB_BACKDOOR AttackType = "JS_EXPLOIT_PROFTPD_BACKDOOR"

	// 外部のjsから攻撃できるよくあるDatabaseのログイン攻撃
	// mssql-default-logins
	// postgres-default-logins
	// redis-default-logins
	// ssh-default-logins
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/javascript/default-logins
	JS_EXPLOIT_DEFAULT_DATABASE_LOGIN AttackType = "JS_EXPLOIT_DEFAULT_DATABASE_LOGIN"

	// プロトコルにおける外部から実行できるjsの攻撃手法
	// mssql-detect
	// oracle-tns-listener
	// samba-detect
	// ssh-auth-methods
	// ref -https://github.com/projectdiscovery/nuclei-texqxmplates/tree/main/javascript/detection
	JS_EXPLOIT_DETECT_MYSQL              AttackType = "JS_EXPLOIT_DETECT_MYSQL"
	JS_EXPLOIT_DETECT_ORACLE_TNS_LISTNER AttackType = "JS_EXPLOIT_DETECT_ORACLE_TNS_LISTENer"
	JS_EXPLOIT_DETECT_SAMBA              AttackType = "JS_EXPLOIT_DETECT_ORACLE_SAMBA"
	JS_EXPLOIT_DETECT_SSH_AUTH           AttackType = "JS_EXPLOIT_DETECT_SSH_AUTH"

	// 各プロトコルでよくあるjsエクスプロイトコード
	// todo: (shinta:enka) beautiful soupで取れる値がわかってきたら、もう少し細かく分ける
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/javascript/enumeration
	JS_EXPLOIT_LDAP  AttackType = "JS_EXPLOIT_LDAP"
	JS_EXPLOIT_PGSQL AttackType = "JS_EXPLOIT_PGSQL"
	JS_EXPLOIT_POP3  AttackType = "JS_EXPLOIT_POP3"
	JS_EXPLOIT_REDIS AttackType = "JS_EXPLOIT_REDIS"
	JS_EXPLOIT_RSYNC AttackType = "JS_EXPLOIT_RSYNC"
	JS_EXPLOIT_SMB   AttackType = "JS_EXPLOIT_SMB"
	JS_EXPLOIT_SSH   AttackType = "JS_EXPLOIT_SSH"

	// cve公開されてる、jsエクスプロイトコード
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/javascript/cves
	JS_EXPLOIT_MYSQL_CVE_2012_2122       AttackType = "JS_EXPLOIT_MYSQL_CVE_2012_2122"
	JS_EXPLOIT_MEMCACHED_CVE_2016_8706   AttackType = "JS_EXPLOIT_MEMCACHED_CVE_2016_8706"
	JS_EXPLOIT_PGSQL_CVE_2019_9193       AttackType = "JS_EXPLOIT_PGSQL_CVE_2019_9193"
	JS_EXPLOIT_OPENSMTPD_CVE_2020_7247   AttackType = "JS_EXPLOIT_OPENSMTPD_CVE_2020_7247"
	JS_EXPLOIT_VMWAREARIA_CVE_2020_34039 AttackType = "JS_EXPLOIT_VMWAREARIA_CVE_2020_34039"
	JS_EXPLOIT_OPENSMTPD_CVE_2023_46604  AttackType = "JS_EXPLOIT_OPENSMTPD_CVE_2023_46604"
	JS_EXPLOIT_OPENSMTPD_CVE_2023_48795  AttackType = "JS_EXPLOIT_OPENSMTPD_CVE_2023_48795"
	JS_EXPLOIT_JENKINS_CVE_2024_23897    AttackType = "JS_EXPLOIT_JENKINS_CVE_2024_23897"

	// server miss configuration
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/javascript/misconfiguration
	JS_EXPLOIT_MISS_CONFIG_SSH   AttackType = "JS_EXPLOIT_MISS_CONFIG_SSH"
	JS_EXPLOIT_MISS_CONFIG_SMB   AttackType = "JS_EXPLOIT_MISS_CONFIG_SMB"
	JS_EXPLOIT_MISS_CONFIG_PGSQL AttackType = "JS_EXPLOIT_MISS_CONFIG_PGSQL"
	JS_EXPLOIT_MISS_CONFIG_MYSQL AttackType = "JS_EXPLOIT_MISS_CONFIG_MYSQL"
	JS_EXPLOIT_MISS_CONFIG_X11   AttackType = "JS_EXPLOIT_MISS_CONFIG_MYSQL"

	// **

	// awsをターゲットにした攻撃手法
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/cloud/aws
	AWS_ACM        AttackType = "AWS_ACM"
	AWS_CLOUDTRAIL AttackType = "AWS_CLOUDTRAIL"
	AWS_CLOUDWATCH AttackType = "AWS_CLOUDWATCH"
	AWS_EC2        AttackType = "AWS_EC2"
	AWS_IAM        AttackType = "AWS_IAM"
	AWS_RDS        AttackType = "AWS_RDS"
	AWS_S3         AttackType = "AWS_S3"
	AWS_SNS        AttackType = "AWS_SNS"
	AWS_VPC        AttackType = "AWS_VPC"
	RT_AWS_VPN     AttackType = "RT_AWS_VPN"

	// k8sをターゲットにした攻撃手法
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/cloud/kubernetes
	CLOUD_K8S_DEPLOYMENTS         AttackType = "CLOUD_K8S_DEPLOYMENTS"
	CLOUD_K8S_NETWORK_POLICY      AttackType = "CLOUD_K8S_NETWORK_POLICY"
	CLOUD_K8S_POD                 AttackType = "CLOUD_K8S_POD"
	CLOUD_K8S_SECURITY_COMPLIANCE AttackType = "CLOUD_K8S_SECURITY_COMPLIANCE"

	// dast攻撃
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/dast/vulnerabilitiesx
	DAST_OAST_POLYGLOT AttackType = "DAST_OAST_POLYGLOT"
	DAST_RUBY_RCE      AttackType = "DAST_RUBY_RCE"

	// dastの公開されているcve
	// todo: (shinta:enka) beautiful soupで取れる値がわかってきたら、もう少し細かく分ける
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/dast/cves
	DAST_CVE AttackType = "DAST_CVE"

	// dnsの攻撃
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/dns
	DNS_AZURE_TAKEOVER             AttackType = "DNS_AZURE_TAKEOVER"
	DNS_BIMI_DETECT                AttackType = "DNS_BIMI_DETECT"
	DNS_CAA_FINGER_PRINT           AttackType = "DNS_CAA_FINGER_PRINT"
	DNS_DETECT_DANGLING_CNAME      AttackType = "DNS_DETECT_DANGLING_CNAME"
	DNS_DMARC_DETECT               AttackType = "DNS_DMARC_DETECT"
	DNS_REBINDING                  AttackType = "DNS_REBINDING"
	DNS_SAAS_SERVICE_DETECTION     AttackType = "DNS_SAAS_SERVICE_DETECTION"
	DNS_WAF_DETECT                 AttackType = "DNS_WAF_DETECT"
	DNS_DNSSEC_DETECTION           AttackType = "DNS_DNSSEC_DETECTION"
	DNS_EC2_DETECTION              AttackType = "DNS_EC2_DETECTION"
	DNS_ELASTIC_BEANSTALK_TAKEOVER AttackType = "DNS_ELASTIC_BEANSTALK_TAKEOVER"
	DNS_MX_FINGERPRINT             AttackType = "DNS_MX_FINGERPRINT"
	DNS_MX_SERVICE_DETECTOR        AttackType = "DNS_MX_SERVICE_DETECTOR"
	DNS_NAMESERVER_FINGERPRINT     AttackType = "DNS_NAMESERVER_FINGERPRINT"
	DNS_PTR_FINGERPRINT            AttackType = "DNS_PTR_FINGERPRINT"
	DNS_SERVFAIL_REFUSED_HOSTS     AttackType = "DNS_SERVFAIL_REFUSED_HOSTS"
	DNS_SOA_DETECT                 AttackType = "DNS_SOA_DETECT"
	DNS_SPF_RECORD_DETECT          AttackType = "DNS_SPF_RECORD_DETECT"
	DNS_SPOOFABLE_SPF_RECORDS_PTR  AttackType = "DNS_SPOOFABLE_SPF_RECORDS_PTR"
	DNS_TXT_FINGERPRINT            AttackType = "DNS_TXT_FINGERPRINT"
	DNS_TXT_SERVICE_DETECT         AttackType = "DNS_TXT_SERVICE_DETECT"
	DNS_WORKSITES_DETECTION        AttackType = "DNS_WORKSITES_DETECTION"

	// sslの攻撃
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/ssl
	SSL_DEPRECATED_TLS               AttackType = "SSL_DEPRECATED_TLS"
	SSL_DETECT_SSL_ISSUER            AttackType = "SSL_DETECT_SSL_ISSUER"
	SSL_EXPIRED_SSL                  AttackType = "SSL_EXPIRED_SSL"
	SSL_INSECURE_CIPHER_SUITE_DETECT AttackType = "SSL_INSECURE_CIPHER_SUITE_DETECT"
	SSL_KUBERNETES_FAKE_CERTIFICATE  AttackType = "SSL_KUBERNETES_FAKE_CERTIFICATE"
	SSL_MISMATCHED_SSL_CERTIFICATE   AttackType = "SSL_MISMATCHED_SSL_CERTIFICATE"
	SSL_REVOKED_SSL_CERTIFICATE      AttackType = "SSL_REVOKED_SSL_CERTIFICATE"
	SSL_SELF_SIGNED_SSL              AttackType = "SSL_SELF_SIGNED_SSL"
	SSL_SSL_DNS_NAMES                AttackType = "SSL_SSL_DNS_NAMES"
	SSL_TLS_VERSION                  AttackType = "SSL_SSL_DNS_NAMES"
	SSL_UNTRUSTED_ROOT_CERTIFICATE   AttackType = "SSL_UNTRUSTED_ROOT_CERTIFICATE"
	SSL_WEAK_CIPHER_SUITES           AttackType = "SSL_WEAK_CIPHER_SUITES"
	SSL_WILDCARD_TLS                 AttackType = "SSL_WILDCARD_TLS"

	// network cve エクスプロイト攻撃
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network

	NETWORK_EXPLOIT_C2_DARKCOMET_TROJAN                       AttackType = "NETWORK_C2_DARKCOMET_TROJAN"
	NETWORK_EXPLOIT_C2_DARKTRACK_RAT_TROJAN                   AttackType = "NETWORK_C2_DARKTRACK_RAT_TROJAN"
	NETWORK_EXPLOIT_C2_ORCUS_RAT_TROJAN                       AttackType = "NETWORK_C2_ORCUS_RAT_TROJAN"
	NETWORK_EXPLOIT_C2_XTREMERAT_TROJAN                       AttackType = "NETWORK_C2_XTREMERAT_TROJAN"
	NETWORK_EXPLOIT_DEPRECATED_SSH_CVE_2001_1473              AttackType = "NETWORK_DEPRECATED_SSH_CVE_2001_1473"
	NETWORK_EXPLOIT_DISTCCDCVE_2004_2687                      AttackType = "NETWORK_DISTCCDCVE_2004_2687"
	NETWORK_EXPLOIT_VSFTPD_CVE_2011_2523                      AttackType = "NETWORK_VSFTPD_CVE_2011_2523"
	NETWORK_EXPLOIT_PROFTPD_CVE_2015_3306                     AttackType = "NETWORK_PROFTPD_CVE_2015_3306"
	NETWORK_EXPLOIT_HP_DATA_PROTECTOR_CVE_2016_2004           AttackType = "NETWORK_HP_DATA_PROTECTOR_CVE_2016_2004"
	NETWORK_EXPLOIT_ORACLE_WEBLOGIC_SERVER_JAVA_CVE_2016_3510 AttackType = "NETWORK_ORACLE_WEBLOGIC_SERVER_JAVA_CVE_2016_3510"
	NETWORK_EXPLOIT_CISCO_IOS_CVE_2017_3881                   AttackType = "NETWORK_CISCO_IOS_CVE_2017_3881"
	NETWORK_EXPLOIT_APACHE_LOG4J_SERVER_CVE_2017_5645         AttackType = "NETWORK_APACHE_LOG4J_SERVER_CVE_2017_5645"
	NETWORK_EXPLOIT_ORACLE_WEBLOGIC_SERVER_CVE_2018_2628      AttackType = "NETWORK_ORACLE_WEBLOGIC_SERVER_CVE_2018_2628"
	NETWORK_EXPLOIT_ORACLE_WEBLOGIC_SERVER_CVE_2018_2893      AttackType = "NETWORK_ORACLE_WEBLOGIC_SERVER_CVE_2018_2893"
	NETWORK_EXPLOIT_APACHE_AIRFLOW_CVE_2020_11981             AttackType = "NETWORK_APACHE_AIRFLOW_CVE_2020_11981"
	NETWORK_EXPLOIT_GHOSTCAT_CVE_2020_1938                    AttackType = "NETWORK_GHOSTCAT_CVE_2020_1938"
	NETWORK_EXPLOIT_APACHE_CASSANDRA_CVE_2021_44521           AttackType = "NETWORK_APACHE_CASSANDRA_CVE_2021_44521"
	NETWORK_EXPLOIT_REDIS_SANDBOX_ESCAPE_CVE_2022_0543        AttackType = "NETWORK_REDIS_SANDBOX_ESCAPE_CVE_2022_0543"
	NETWORK_EXPLOIT_COUCHDB_ERLANG_CVE_2022_24706             AttackType = "NETWORK_COUCHDB_ERLANG_CVE_2022_24706"
	NETWORK_EXPLOIT_MUHTTPD_CVE_2022_31793                    AttackType = "NETWORK_MUHTTPD_CVE_2022_31793"
	NETWORK_EXPLOIT_ROCKETMQ_CVE_2023_33246                   AttackType = "NETWORK_EXPLOIT_ROCKETMQ_CVE_2023_33246"

	// networkにおけるdefault login攻撃
	// todo: (shinta:enka) beautiful soupで取れる値がわかってきたら、もう少し細かく分ける
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network/default-login
	NETWORK_DEFAULT_LOGIN AttackType = "NETWORK_DEFAULT_LOGIN"

	// network detect 攻撃
	// todo: (shinta:enka) beautiful soupで取れる値がわかってきたら、もう少し細かく分ける
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network/detection
	NETWORK_EXPLOIT_DETECT AttackType = "NETWORK_EXPLOIT_DETECT"

	// network越しに実行可能な攻撃
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network/enumeration
	NETWORK_EXPLOIT_SMTP                  AttackType = "NETWORK_SMTP"
	NETWORK_EXPLOIT_BEANSTALK_SERVICE     AttackType = "NETWORK_BEANSTALK_SERVICE"
	NETWORK_EXPLOIT_KAFKA_TOPIC           AttackType = "NETWORK_KAFKA_TOPIC"
	NETWORK_EXPLOIT_MONGODB_INFO_ENUM     AttackType = "NETWORK_MONGODB_INFO_ENUM"
	NETWORK_EXPLOIT_NIAGARA_FOX_INFO_ENUM AttackType = "NETWORK_NIAGARA_FOX_INFO_ENUM"
	NETWORK_EXPLOIT_PGSQL_USER_ENUM       AttackType = "NETWORK_NIAGARA_FOX_INFO_ENUM"

	// network越しに可能なexposure攻撃一覧
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network/exposures
	NETWORK_CISCO_SMI_EXPOSURE AttackType = "NETWORK_CISCO_SMI_EXPOSURE"
	NETWORK_EXPOSED_ADB        AttackType = "NETWORK_EXPOSED_ADB"
	NETWORK_EXPOSED_DOCKERD    AttackType = "NETWORK_EXPOSED_DOCKERD"
	NETWORK_EXPOSED_REDIS      AttackType = "NETWORK_EXPOSED_REDIS"
	NETWORK_EXPOSED_ZOOKEEPER  AttackType = "NETWORK_EXPOSED_ZOOKEEPER"

	// network越しに設置可能なハニーポット攻撃一覧
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network/honeypot
	NETWORK_HONEYPOT_ADBHONEY_HONEYPOT_CNXN_DETECT  AttackType = "NETWORK_HONEYPOT_ADBHONEY_HONEYPOT_CNXN_DETECT"
	NETWORK_HONEYPOT_ADBHONEY_HONEYPOT_SHELL_DETECT AttackType = "NETWORK_HONEYPOT_ADBHONEY_HONEYPOT_SHELL_DETECT"
	NETWORK_HONEYPOT_CONPOT_SIEMENS_HONEYPOT_DETECT AttackType = "NETWORK_HONEYPOT_CONPOT_SIEMENS_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_COWRIE_SSH_HONEYPOT_DETECT     AttackType = "NETWORK_HONEYPOT_COWRIE_SSH_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_DIONAEA_FTP_HONEYPOT_DETECT    AttackType = "NETWORK_HONEYPOT_DIONAEA_FTP_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_DIONAEA_MQTT_HONEYPOT_DETECT   AttackType = "NETWORK_HONEYPOT_DIONAEA_MQTT_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_DIONAEA_MYSQL_HONEYPOT_DETECT  AttackType = "NETWORK_HONEYPOT_DIONAEA_MYSQL_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_DIONAEA_SMB_HONEYPOT_DETECT    AttackType = "NETWORK_HONEYPOT_DIONAEA_SMB_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_GASPOT_HONEYPOT_DETECT         AttackType = "NETWORK_HONEYPOT_GASPOT_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_MAILONEY_HONEYPOT_DETECT       AttackType = "NETWORK_HONEYPOT_MAILONEY_HONEYPOT_DETECT"
	NETWORK_HONEYPOT_REDIS_HONEYPOT_DETECT          AttackType = "NETWORK_HONEYPOT_REDIS_HONEYPOT_DETECT"

	// network越しにjarm実行可能な攻撃一覧
	// JARM（JA3/SSL Fingerprinting）とは、特定のサーバーが使用するTLS（Transport Layer Security）ハンドシェイクの特性を収集し、そのサーバーの識別情報（指紋）を生成する技術
	// サーバーが使用する暗号化設定や証明書、TLSバージョンなどの情報を元に生成されます
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network/honeypot
	NETWORK_JARM_COBALT_STRIKE_C2_JARM AttackType = "NETWORK_JARM_COBALT_STRIKE_C2_JARM"
	NETWORK_JARM_COVENANT_C2_JARM      AttackType = "NETWORK_JARM_COVENANT_C2_JARM"
	NETWORK_JARM_DEIMOS_C2_JARM        AttackType = "NETWORK_JARM_DEIMOS_C2_JARM"
	NETWORK_JARM_EVILGINX2_JARM        AttackType = "NETWORK_JARM_EVILGINX2_JARM"
	NETWORK_JARM_GENERIC_C2_JARM       AttackType = "NETWORK_JARM_GENERIC_C2_JARM"
	NETWORK_JARM_GRAT2_C2_JARM         AttackType = "NETWORK_JARM_GRAT2_C2_JARM"
	NETWORK_JARM_HAVOC_C2_JARM         AttackType = "NETWORK_JARM_HAVOC_C2_JARM"
	NETWORK_JARM_MAC_C2_JARM           AttackType = "NETWORK_JARM_MAC_C2_JARM"
	NETWORK_JARM_MACSHELL_C2_JARM      AttackType = "NETWORK_JARM_MACSHELL_C2_JARM"
	NETWORK_JARM_MERLIN_C2_JARM        AttackType = "NETWORK_JARM_MERLIN_C2_JARM"
	NETWORK_JARM_METASPLOIT_C2_JARM    AttackType = "NETWORK_JARM_METASPLOIT_C2_JARM"
	NETWORK_JARM_MYTHIC_C2_JARM        AttackType = "NETWORK_JARM_MYTHIC_C2_JARM"
	NETWORK_JARM_POSH_C2_JARM          AttackType = "NETWORK_JARM_POSH_C2_JARM"
	NETWORK_JARM_SHAD0W_C2_JARM        AttackType = "NETWORK_JARM_SHAD0W_C2_JARM"
	NETWORK_JARM_SILENTTRINITY_C2_JARM AttackType = "NETWORK_JARM_SILENTTRINITY_C2_JARM"
	NETWORK_JARM_SLIVER_C2_JARM        AttackType = "NETWORK_JARM_SLIVER_C2_JARM"

	// network越しに検出可能なmiss config
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/network/misconfig
	NETWORK_MISS_CONFIG                               AttackType = "NETWORK_MISS_CONFIG"
	NETWORK_MISS_CONFIG_APACHE_DUBBO_UNAUTH           AttackType = "NETWORK_MISS_CONFIG_APACHE_DUBBO_UNAUTH"
	NETWORK_MISS_CONFIG_APACHE_ROCKETMQ_BROKER_UNAUTH AttackType = "NETWORK_MISS_CONFIG_APACHE_ROCKETMQ_BROKER_UNAUTH"
	NETWORK_MISS_CONFIG_CLAMAV_UNAUTH                 AttackType = "NETWORK_MISS_CONFIG_CLAMAV_UNAUTH"
	NETWORK_MISS_CONFIG_CLICKHOUSE_UNAUTH             AttackType = "NETWORK_MISS_CONFIG_CLICKHOUSE_UNAUTH"
	NETWORK_MISS_CONFIG_ERLANG_DAEMON                 AttackType = "NETWORK_MISS_CONFIG_ERLANG_DAEMON"
	NETWORK_MISS_CONFIG_GANGLIA_XML_GRID_MONITOR      AttackType = "NETWORK_MISS_CONFIG_GANGLIA_XML_GRID_MONITOR"
	NETWORK_MISS_CONFIG_MEMCACHED_STATS               AttackType = "NETWORK_MISS_CONFIG_MEMCACHED_STATS"
	NETWORK_MISS_CONFIG_MONGODB_UNAUTH                AttackType = "NETWORK_MISS_CONFIG_MONGODB_UNAUTH"
	NETWORK_MISS_CONFIG_MYSQL_NATIVE_PASSWORD         AttackType = "NETWORK_MISS_CONFIG_MYSQL_NATIVE_PASSWORD"
	NETWORK_MISS_CONFIG_PRINTERS_INFO_LEAK            AttackType = "NETWORK_MISS_CONFIG_PRINTERS_INFO_LEAK"
	NETWORK_MISS_CONFIG_SAP_ROUTER_INFO_LEAK          AttackType = "NETWORK_MISS_CONFIG_SAP_ROUTER_INFO_LEAK"
	NETWORK_MISS_CONFIG_TIDB_NATIVE_PASSWORD          AttackType = "NETWORK_MISS_CONFIG_TIDB_NATIVE_PASSWORD"
	NETWORK_MISS_CONFIG_TIDB_UNAUTH                   AttackType = "NETWORK_MISS_CONFIG_TIDB_UNAUTH"
	NETWORK_MISS_CONFIG_UNAUTH_PSQL                   AttackType = "NETWORK_MISS_CONFIG_UNAUTH_PSQL"

	// network越しに検出可能なvulnerabilitiese
	// ref - https://github.com/projectdiscovery/nuclei-templates/blob/main/network/vulnerabilities/clockwatch-enterprise-rce.yaml
	NETWORK_VULNERABILITIES_CLOCKWATCH_ENTERPRISE_RCE AttackType = "NETWORK_JARM_SLIVER_C2_JARM"

	// http越しに攻撃可能な攻撃
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/http

	// file攻撃, e.g zip爆弾など
	// ref - https://github.com/projectdiscovery/nuclei-templates/tree/main/file
)

// ClassficatioAttackTypesはlove machine mlで使用できる、AttackTypeを
// class化した、AttackTypeの集合値
// マッチする条件においてパターン化された攻撃手法のまとめり
// 将来的にwhite-hat-moduleで最適化されるべき箇所である。
type ClassficatioAttackTypes []AttackType
