# error-reporting-sample

# CloudRunにデプロイ
https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service?hl=ja

```
gcloud config set project nekoshita-error-reporting-test
asdf install
go mod init github.com/nekoshita/error-reporting-sample
gcloud run deploy error-reporting-sample --region asia-southeast1 --project nekoshita-error-reporting-test
curl https://error-reporting-sample-3ii3ii6eda-as.a.run.app
```

# Error Reporting API
```
gcloud beta error-reporting events report --project nekoshita-error-reporting-test --service=SERVICE_NAME --message='panic: hoge

goroutine 1 [running]:
main.main()
        /github.com/nekoshita/error-reporting-sample/main.go:12 +0x2c
exit status 2'
```

```
gcloud beta error-reporting events report --project nekoshita-error-reporting-test --service=SERVICE_NAME --message='panic: hoge
goroutine 1 [running]:
main.main()
        /github.com/nekoshita/error-reporting-sample/main.go:12 +0x2c
exit status 2'
```

```
gcloud beta error-reporting events report --project nekoshita-error-reporting-test --service=SERVICE_NAME --message='java.lang.CloneNotSupportedException
at sample.24.SmapleUserList.Loging(SampleUserList.java:27)'
```
