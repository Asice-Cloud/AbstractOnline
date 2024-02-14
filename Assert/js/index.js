alert("Hello, World!")


window.onload=function f(){
    var myDate=new Date()
    document.getElementById("currentTime").innerText=myDate.getDate()+"" +
        "-"+(myDate.getMonth()+1)+"-"+myDate.getFullYear()+" "+myDate.getHours()+":"+myDate.getMinutes()+":"+myDate.getSeconds();
}
function loadPage() {
    var targetURL=document.getElementById("url").value;
    console.log(targetURL);
    document.getElementById("iframePosition").src=targetURL;
}