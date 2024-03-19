### Install
```
go install github.com/cosmtrek/air@latest
go get .
go mod tidy
npm install
```

### Development
```
npm run dev
```

### Building
```
npm run tailwind
go build .
```


```
javascript:async function savePage(){let t=await (await fetch("http://localhost:8090/add",{method:"POST",body:JSON.stringify({url:window.location.href,html:document.documentElement.innerHTML})})).text();window.open("http://localhost:8090/edit/"+t)}savePage();
```