Install
```
go install github.com/cosmtrek/air@latest
npm install
```

Dev
```
npm run dev
```


go install github.com/cosmtrek/air@latest

javascript:async function savePage(){let t=await (await fetch("http://localhost:8090/add",{method:"POST",body:JSON.stringify({url:window.location.href,html:document.documentElement.innerHTML})})).text();window.open("http://localhost:8090/edit/"+t)}savePage();