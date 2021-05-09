baseURL = document.getElementById("baseURL").innerHTML
document.getElementById("uploadcmd").innerHTML = "curl -F 'file=@uploadthis.txt' " + baseURL

function submitForm() {
    var f = document.getElementsByName('upload-form')[0]
    f.submit()
    f.reset()
}

function copyText() {
    const tmp = document.createElement('textarea')
    tmp.value = (document.getElementById("uploadcmd")).innerHTML
    document.body.appendChild(tmp)
    tmp.select()
    document.execCommand('copy')
    document.body.removeChild(tmp)
}

function getCellIdx(id) {
    if (id == "size" || id == "modtime") {
        return document.getElementById("th_"+id).cellIndex + 1
    }

    return document.getElementById("th_"+id).cellIndex
}

function sortTable(id, type) {
    var rows, i, x, y, shouldSwitch
    table = document.getElementById("filetable")
    sortorder = document.getElementById("sort_" + id)
    cellIdx = getCellIdx(id)
    switching = true
    while (switching) {
        switching = false
        rows = table.rows
        for (i = 1; i < (rows.length - 1); i++) {
            shouldSwitch = false
            x = rows[i].getElementsByTagName("TD")[cellIdx]
            y = rows[i + 1].getElementsByTagName("TD")[cellIdx]

            if (type == "num") {
                if (Number(x.innerHTML) > Number(y.innerHTML) && sortorder.innerHTML > 0) {
                    shouldSwitch = true
                    break
                }
                
                if (Number(x.innerHTML) < Number(y.innerHTML) && sortorder.innerHTML < 0) {
                    shouldSwitch = true
                    break
                }
            } else if (type == "date"){
                x = new Date(x.innerHTML)
                y = new Date(y.innerHTML)
                if (x > y && sortorder.innerHTML > 0) {
                    shouldSwitch = true
                    break
                }
                
                if (x < y && sortorder.innerHTML < 0) {
                    shouldSwitch = true
                    break
                }
            } else {
                if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase() && sortorder.innerHTML > 0) {
                    shouldSwitch = true
                    break
                }

                if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase() && sortorder.innerHTML < 0) {
                    shouldSwitch = true
                    break
                }
            }
        }
        if (shouldSwitch) {
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i])
            switching = true
        }
    }
    sortorder.innerHTML = Number(sortorder.innerHTML) * -1
}

function showQR(url) {
    qrdiv = document.getElementById("qrcode")
    qrcode = new QRCode(qrdiv, {width: 256, height: 256, margin: "auto"})
    
    var modal = document.getElementById("modalPopUp")
    var captionText = document.getElementById("caption")

    modal.style.display = "block"
    captionText.innerHTML = baseURL + url

    var span = document.getElementsByClassName("close")[0]
    span.onclick = function() { 
        modal.style.display = "none"
        qrdiv.innerHTML = ""
        qrcode.clear()
    }

    qrcode.makeCode(baseURL + url)
}