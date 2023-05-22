//import * as params from '@params';

var resList = document.getElementById('searchResults');
var sInput = document.getElementById('searchInput');
var first, last, current_elem = null
var resultsAvailable = false;

// http request

var http = {};
http.quest = function (option, callback) {
    var url = option.url;
    var method = option.method;
    var data = option.data;
    var timeout = option.timeout || 0;
    var xhr = new XMLHttpRequest();
    (timeout > 0) && (xhr.timeout = timeout);
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4) {
            if (xhr.status >= 200 && xhr.status < 400) {
                var result = xhr.responseText;
                try { result = JSON.parse(xhr.responseText); } catch (e) { }
                callback && callback(null, result);
            } else {
                callback && callback('status: ' + xhr.status);
            }
        }
    }.bind(this);
    xhr.open(method, url, true);
    if (typeof data === 'object') {
        try {
            data = JSON.stringify(data);
        } catch (e) { }
    }
    xhr.send(data);
    xhr.ontimeout = function () {
        callback && callback('timeout');
        console.log('%c连%c接%c超%c时', 'color:red', 'color:orange', 'color:purple', 'color:green');
    };
};
http.get = function (url, callback) {
    var option = url.url ? url : { url: url };
    option.method = 'get';
    this.quest(option, callback);
};


function search(query){

    http.get('../api/search?q='+query, function (err, result) {
        if (err!==null){
            return
        }
        var results = []
        if (result.code!==0){
            return
        }
        if (result.data.results===undefined){
            return
        }
        results=result.data.results
        if (results.length !== 0) {
            // build our html if result exists
            let resultSet = ''; // our results bucket

            for (let item in results) {
                resultSet += `<li class="post-entry"><header class="entry-header">${results[item].title}&nbsp;»</header>` +
                    `<a href="${results[item].link}" aria-label="${results[item].title}"></a></li>`
            }

            resList.innerHTML = resultSet;
            resultsAvailable = true;
            first = resList.firstChild;
            last = resList.lastChild;
        } else {
            resultsAvailable = false;
            resList.innerHTML = '';
        }
    });
}



function activeToggle(ae) {
    document.querySelectorAll('.focus').forEach(function (element) {
        // rm focus class
        element.classList.remove("focus")
    });
    if (ae) {
        ae.focus()
        document.activeElement = current_elem = ae;
        ae.parentElement.classList.add("focus")
    } else {
        document.activeElement.parentElement.classList.add("focus")
    }
}

function reset() {
    resultsAvailable = false;
    resList.innerHTML = sInput.value = ''; // clear inputbox and searchResults
    sInput.focus(); // shift focus to input box
}


sInput.onkeyup = function (e) {
    search(this.value.trim()) // the actual query being run using fuse.js
}

sInput.addEventListener('search', function (e) {
    // clicked on x
    if (!this.value) reset()
})

// kb bindings
document.onkeydown = function (e) {
    let key = e.key;
    var ae = document.activeElement;

    let inbox = document.getElementById("searchbox").contains(ae)

    if (ae === sInput) {
        var elements = document.getElementsByClassName('focus');
        while (elements.length > 0) {
            elements[0].classList.remove('focus');
        }
    } else if (current_elem) ae = current_elem;

    if (key === "Escape") {
        reset()
    } else if (!resultsAvailable || !inbox) {
        return
    } else if (key === "ArrowDown") {
        e.preventDefault();
        if (ae == sInput) {
            // if the currently focused element is the search input, focus the <a> of first <li>
            activeToggle(resList.firstChild.lastChild);
        } else if (ae.parentElement != last) {
            // if the currently focused element's parent is last, do nothing
            // otherwise select the next search result
            activeToggle(ae.parentElement.nextSibling.lastChild);
        }
    } else if (key === "ArrowUp") {
        e.preventDefault();
        if (ae.parentElement == first) {
            // if the currently focused element is first item, go to input box
            activeToggle(sInput);
        } else if (ae != sInput) {
            // if the currently focused element is input box, do nothing
            // otherwise select the previous search result
            activeToggle(ae.parentElement.previousSibling.lastChild);
        }
    } else if (key === "ArrowRight") {
        ae.click(); // click on active link
    }
}
