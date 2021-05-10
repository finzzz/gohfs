baseURL = document.getElementById("baseURL").innerHTML
zipPath = document.getElementById("zipPath")
document.getElementById("up_curl").innerHTML = "curl -F 'file=@uploadthis.txt' " + baseURL

renderItems("linksvg", "static/icons/link.svg")
renderItems("zipsvg", "static/icons/zip.svg")
renderItems("qrsvg", "static/icons/qr.svg")
renderItems("hashsvg", "static/icons/hash.svg")
renderItems("termsvg", "static/icons/term.svg")

setZipPath()

function setZipPath() {
    var link = document.getElementsByClassName("ziplink")

    Array.prototype.forEach.call(link, function(slide, index) {
        link.item(index).href = zipPath.innerHTML + "/" + link.item(index).name
    });
}

function renderItems(cls, url) {
    var webPath = document.getElementById("webPath")
    var link = document.getElementsByClassName(cls)

    Array.prototype.forEach.call(link, function(slide, index) {
        link.item(index).src = webPath.innerHTML + "/" + url
    });
}

function submitForm() {
    var f = document.getElementsByName('upload-form')[0]
    f.submit()
    f.reset()
}

function copyTextAsURL(text) {
    const tmp = document.createElement('textarea')
    tmp.value = baseURL + "/" + text
    document.body.appendChild(tmp)
    tmp.select()
    document.execCommand('copy')
    document.body.removeChild(tmp)
}

function copyTextById(id) {
    const tmp = document.createElement('textarea')
    tmp.value = (document.getElementById(id)).innerHTML
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

function showModal(id) {
    var modal = document.getElementById(id)
    modal.style.display = "flex"
}

function hideModal(id) {
    var modal = document.getElementById(id)
    modal.style.display = "none"

    document.getElementById("qrcode").innerHTML = ""
    document.getElementById("zipqrcode").innerHTML = ""
}

function showQR(url) {
    qrdiv = document.getElementById("qrcode")
    qrcode = new QRCode(qrdiv, {width: 256, height: 256, margin: "auto"})
    document.getElementById("caption").innerHTML = "RAW : " + baseURL + url

    // zip
    zipqrdiv = document.getElementById("zipqrcode")
    zipqrcode = new QRCode(zipqrdiv, {width: 256, height: 256, margin: "auto"})
    document.getElementById("zipcaption").innerHTML = "ZIP : " + baseURL + zipPath.innerHTML + url

    showModal("QRModal")

    qrcode.makeCode(baseURL + url)
    zipqrcode.makeCode(baseURL + zipPath.innerHTML + url)
}