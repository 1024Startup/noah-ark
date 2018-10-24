'use strict';

/**
 *
 * vue工具函数
 *
 **/

/**
 * @Author: jeffery
 * @Date:   2016-12-07T10:27:59+08:00
 * @Email:  jeffery@influx.io
 * @Filename: custom.js
 * @Last modified by:   jeffery
 * @Last modified time: 2017-08-07T16:52:04+08:00
 * @Copyright: Copyright 2016 Jeffery. All rights reserved.
 */

var qatrix = function qatrix() {};

/**
 * 对象安全访问
 * @param  {[type]} p  访问路径数组
 * @param  {[type]} o  访问对象
 * @return {[type]}    [description] example:var obj = {a:{b:[0]}}   qatrix.get(['a','b',0],obj) will return 0
 */
qatrix['get'] = function (p, o) {
    return p.reduce(function (xs, x) {
        return xs && xs[x] || xs && xs[x] === 0 ? xs[x] : null;
    }, o);
};

/**
 * querystring handler
 * @type {Object}
 */
qatrix['qs'] = {
    /**
     * Turn the given `obj` into a query string
     *
     * @param {Object} obj
     * @return {String}
     * @api public
     */
    stringify: function stringify(obj, prefix) {
        var self = this;
        if ('[object Array]' == toString.call(obj)) {
            return self.stringifyArray(obj, prefix);
        } else if ('[object Object]' == toString.call(obj)) {
            return self.stringifyObject(obj, prefix);
        } else if ('string' == typeof obj) {
            return self.stringifyString(obj, prefix);
        } else {
            return prefix + '=' + String(obj);
        }
    },
    /**
     * Stringify the given `str`.
     *
     * @param {String} str
     * @param {String} prefix
     * @return {String}
     * @api private
     */
    stringifyString: function stringifyString(str, prefix) {
        if (!prefix) throw new TypeError('stringify expects an object');
        return prefix + '=' + str;
    },
    /**
     * Stringify the given `arr`.
     *
     * @param {Array} arr
     * @param {String} prefix
     * @return {String}
     * @api private
     */
    stringifyArray: function stringifyArray(arr, prefix) {
        var self = this;
        var ret = [];
        if (!prefix) throw new TypeError('stringify expects an object');
        ret.push(JSON.stringify(arr));
        return prefix + "=" + ret.join('&');
    },
    /**
     * Stringify the given `obj`.
     *
     * @param {Object} obj
     * @param {String} prefix
     * @return {String}
     * @api private
     */
    stringifyObject: function stringifyObject(obj, prefix) {
        var self = this;
        var ret = [],
            keys = Object.keys(obj),
            key = void 0;

        for (var i = 0, len = keys.length; i < len; ++i) {
            key = keys[i];
            if ('' == key) continue;
            if (null == obj[key]) {
                ret.push(key + '=');
            } else {
                ret.push(self.stringify('string' == typeof obj[key] ? obj[key] : JSON.stringify(obj[key]), key));
            }
        }
        return ret.join('&');
    },
    /**
     * 获取querystring中的参数内容
     * @param  {String} querystring 查询参数
     * @param  {String} name 名称
     * @return {String}      内容
     */
    getQs: function getQs(url, name) {
        var result = url.match(new RegExp("[\?\&]" + name + "=([^\&]+)", "i"));
        return result == null || result.length < 1 ? '' : result[1];
    },
    /**
     * 解析URL链接中的query string并返回一个对象
     * @param  {[type]} url [description]
     * @return {[type]}     [description]
     */
    parseQueryString: function parseQueryString(url) {
        var self = this;
        var searchObject = {},
            queries = void 0,
            querystr = void 0;
        querystr = url.split('?')[1];
        if (!querystr) return {};
        return self.parseQs(querystr);
    },
    /**
     * 解析URL链接中的query string并返回一个对象
     * @param  {[type]} url [description]
     * @return {[type]}     [description]
     */
    parseQs: function parseQs(querystr) {
        var searchObject = {},
            queries = void 0;
        if (!querystr) return {};
        queries = querystr.split('&');
        for (var i = 0; i < queries.length; i++) {
            var split = queries[i].split('=');
            if (split.length < 2) continue;

            searchObject[split[0]] = split[1];
        }
        return searchObject;
    }
};
/**
 * 文件上传插件
 * @param {Object} data 提交的参数
 * @param {String} url 上传的地址
 * @param {Function} progress 回调函数
 * qatrix.upload({
 *  progress:function(percent){
 *    $("#progress").css({"width":percent});
 *  },
 *  data:{
 *  "files[]":files,//文件数组
 *  "file":file,//单个文件
 *  "others":others,//其他参数
 *  },
 *  async:false,//默认为true
 *  url:'http://baidu.com'
 *}).then(function(res){});
 */
qatrix.upload = function (opts) {
    return new Promise(function (resolve, reject) {
        // var upload_url = location.protocol+"//"+opts.url;
        var upload_url = opts.url;
        var _async = opts.async ? opts.async : true;
        if (!(opts instanceof Object)) throw new Error("Arguments is not object!");
        if (!opts.hasOwnProperty("data")) throw new Error("Data is Null!");

        var xhr = new XMLHttpRequest();
        xhr.open('POST', upload_url, _async);
        // xhr.setRequestHeader("Authorization", qatrix.Cookies.get('token'));
        for(var key in opts.headers){
            xhr.setRequestHeader(key,opts.headers[key])
        }

        var formData = new FormData(),
            startDate;

        //上传文件参数
        for (var i in opts.data) {
            //支持多文件上传
            if (Object.prototype.toString.call(opts.data[i]) === '[object FileList]') {
                for (var j = 0, len = opts.data[i].length; j < len; j++) {
                    formData.append(i, opts.data[i][j]);
                }
            } else {
                formData.append(i, opts.data[i]);
            }
        };
        //间隔时间
        var taking;

        //上传进度事件
        xhr.upload.addEventListener("progress", function (event) {

            if (event.lengthComputable) {
                //计算上传速率:已经完成文件大小/1024/(当前时间－开始上传时间)
                //(loaded)/1024/(nowDate - startDate)
                var nowDate = new Date().getTime();
                taking = nowDate - startDate;
                var x = event.loaded / 1024;
                var y = taking / 1000;
                var uploadSpeed = x / y;
                var formatSpeed;
                if (uploadSpeed > 1024) {
                    formatSpeed = (uploadSpeed / 1024).toFixed(2) + "Mb\/s";
                } else {
                    formatSpeed = uploadSpeed.toFixed(2) + "Kb\/s";
                }
                //event.loaded已经完成上传的大小
                var percentComplete = Math.round(event.loaded * 100 / event.total);
                // progressInnerDiv.style.width = percentComplete + "%";
                // progressText.innerHTML = percentComplete + "%";
                opts.progress(percentComplete + "%");
            }
        }, false);

        xhr.onreadystatechange = function (response) {
            if (xhr.readyState == 4 && xhr.status == 201 && xhr.responseText) {
                //上传文件成功只返回 201
                try {
                    //let url = decodeURIComponent(encodeURIComponent(xhr.responseText))
                    //console.log('url',decodeURIComponent(encodeURIComponent(xhr.responseText)))
                    var url = xhr.responseText.replace(/\\|"/g, ''); //坑太深了
                    resolve(url);
                } catch (e) {
                    reject(e);
                }
            } else if (xhr.readyState == 4 && xhr.status != 201 && xhr.responseText) {
                reject(xhr);
            }
        };
        startDate = parseInt(+new Date() / 1000);
        xhr.send(formData);
    });
};
//获取url链?后面的命名参数值
qatrix.getQs = function (name) {
    var result = location.search.match(new RegExp("[\?\&]" + name + "=([^\&]+)", "i"));
    if (result == null || result.length < 1) {
        return "";
    }
    return result[1];
};
//删除url链?后面的命名参数值
qatrix.delQs = function (name) {
    return location.search.replace(new RegExp("[\?\&]" + name + "=([^\&]+)", "i"), '');
};
/*
 * 转换object为URL
 */

qatrix.objToUrl = function(obj) {
        return Object.keys(obj).map(function(key) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]);
        }).join('&');
    }
/*
 * 获取时间
 * qatrix.formatDate('yyyy-mm-dd hh:ii:ss',time);
 */
qatrix.formatDate = function (format_str, UNIX_timestamp) {
    var a = new Date(UNIX_timestamp * 1000);

    var o = {
        "y+": a.getFullYear(),
        "m+": a.getMonth() + 1, //月份
        "d+": a.getDate(), //日
        "h+": a.getHours(), //小时
        "i+": a.getMinutes(), //分
        "s+": a.getSeconds() //秒
    };
    for (var k in o) {
        if (new RegExp("(" + k + ")", "i").test(format_str)) {
            format_str = format_str.replace(RegExp.$1, o[k].toString().length < 2 ? '0' + o[k] : o[k]);
        }
    }
    return format_str;
};
/*
 * 转换url为Object
 */
qatrix.urlToObj = function () {
    var search = location.search.substring(1);
    var result_obj = search ? JSON.parse('{"' + search.replace(/&/g, '","').replace(/=/g, '":"') + '"}', function (key, value) {
        return key === "" ? value : decodeURIComponent(value);
    }) : {};
    return result_obj;
};

/**
 * webp图片格式检测
 *  @param  {Int} thumbnail 默认750,提供700,400
 */
qatrix.webpJS = {
    ua: navigator.userAgent.toLowerCase(),
    isSupportWebp: false,
    thumbnail: [400, 750], //支持尺寸400x400,750x750
    init: function init() {
        var self = this;
        self.supportWebpTest();
    },
    handler: function handler(url, thumbnail) {
        var self = this;
        var src = url;

        if (self.thumbnail.indexOf(thumbnail) == -1) {
            thumbnail = 750;
        }
        if (self.isSupportWebp && self.imgUrlMatch(src)) {
            var _thumbnail_webp = '!webp' + thumbnail;
            src = src.replace(/\/\/res\./, '//res2.');
            src += src.indexOf(_thumbnail_webp) > -1 ? "" : _thumbnail_webp;
        } else if (self.imgUrlMatch(src)) {
            var _thumbnail_jpg = '!jpg' + thumbnail;
            src = src.replace(/\/\/res\./, '//res2.');
            src += src.indexOf(_thumbnail_jpg) > -1 ? "" : _thumbnail_jpg;
        }
        return src;
    },
    imgUrlMatch: function imgUrlMatch(url) {
        var self = this;
        var regex = new RegExp('.(jpeg|png|jpg)$', 'i');
        var result = url.match(regex);
        return result ? true : false;
    },
    supportWebpTest: function supportWebpTest() {
        var self = this;
        var storage = window.localStorage;
        if (storage) {
            var _supportwebp = storage.getItem('supportwebp');
            if (_supportwebp == null) {
                _supportwebp = self.supportBrowser();
                _supportwebp && storage.setItem('supportwebp', _supportwebp); //检测是否是支持的浏览器
                !_supportwebp && self.loadWebp(function (flag) {
                    self.isSupportWebp = _supportwebp = flag;
                    storage.setItem('supportwebp', _supportwebp);
                });
            }
            self.isSupportWebp = _supportwebp == "true" ? true : false;
        } else {
            _supportwebp = self.supportBrowser();
            !_supportwebp && self.loadWebp(function (flag) {
                self.isSupportWebp = _supportwebp = flag;
            });
            self.isSupportWebp = _supportwebp;
        }
        return self.isSupportWebp;
    },
    loadWebp: function loadWebp(cb) {
        var self = this;
        var img = new Image();
        img.src = "data:image/webp;base64,UklGRjoAAABXRUJQVlA4IC4AAACyAgCdASoCAAIALmk0mk0iIiIiIgBoSygABc6WWgAA/veff/0PP8bA//LwYAAA";
        img.onload = img.onerror = function () {
            cb(img.width === 2 && img.height === 2 ? true : false);
        };
    },
    supportBrowser: function supportBrowser() {
        var self = this;
        self.isSupportWebp = window.chrome || window.opera ? true : false;
        return self.isSupportWebp;
    }
};
qatrix.webpJS.init();

/**
 * 多个数组组合
 * @param {...} array1,array2,....
 * @return {Array} multiple array
 */
qatrix.combination = function (arrays) {
    var r = [],
        args = arrays,
        max = args.length - 1;

    function helper(arr, i) {
        for (var j = 0, l = args[i].length; j < l; j++) {
            var a = arr.slice(0); // clone arr
            a.push(args[i][j]);
            if (i == max) r.push(a);else helper(a, i + 1);
        }
    }
    helper([], 0);
    return r;
};
/**
 * localStorage 简单封装
 * @Function set @param key,value @return true
 * @Function get @param key @return {mixed} obj
 * @Function getAll @param {} @return {mixed} obj
 * @Function remove @param {} @return {mixed} obj
 */
qatrix.storage = {
    store: window.localStorage,
    set: function set(key, value) {
        var _self = this;
        _self.store.setItem(key, JSON.stringify(value));
        return true;
    },
    get: function get(key) {
        var _self = this;
        var value = _self.store.getItem(key);
        return value ? JSON.parse(value) : null; //JSON.parse 不能处理`""`所以需要判断是否为空
    },
    getAll: function getAll() {
        var _self = this;
        var _keys = _self.store;
        var obj = {};
        for (var i in _keys) {
            obj[i] = _keys[i] ? JSON.parse(_keys[i]) : null;
        }
        return obj;
    },
    remove: function remove(key) {
        var _self = this;
        if (key) _self.store.removeItem(key);
        return true;
    },
    removeAll: function removeAll() {
        var _self = this;
        var _keys = _self.store;
        for (var i in _keys) {
            _self.remove(i);
        }
        return true;
    }
};
/**
 * sessionStorage 简单封装
 * @Function set @param key,value @return true
 * @Function get @param key @return {mixed} obj
 * @Function getAll @param {} @return {mixed} obj
 * @Function remove @param {} @return {mixed} obj
 */
qatrix.sessionStorage = {
    store: window.sessionStorage,
    set: function set(key, value) {
        var _self = this;
        _self.store.setItem(key, JSON.stringify(value));
        return true;
    },
    get: function get(key) {
        var _self = this;
        var value = _self.store.getItem(key);
        return value ? JSON.parse(value) : null; //JSON.parse 不能处理`""`所以需要判断是否为空
    },
    getAll: function getAll() {
        var _self = this;
        var _keys = _self.store;
        var obj = {};
        for (var i in _keys) {
            obj[i] = _keys[i] ? JSON.parse(_keys[i]) : null;
        }
        return obj;
    },
    remove: function remove(key) {
        var _self = this;
        if (key) _self.store.removeItem(key);
        return true;
    },
    removeAll: function removeAll() {
        var _self = this;
        var _keys = _self.store;
        for (var i in _keys) {
            _self.remove(i);
        }
        return true;
    }
};
/**
 * hash string
 * @param {String} str
 * @return {String} hash
 */
qatrix.hash = function (str) {
    var hash = 0,
        str = str ? str : "",
        i,
        char;
    if (str.length === 0) return hash;
    for (i = 0; i < str.length; i++) {
        char = str.charCodeAt(i);
        hash = (hash << 5) - hash + char;
        hash |= 0;
    }
    return Math.abs(hash);
};

/**
 * md5方法封装。
 * JavaScript MD5
 * https://github.com/blueimp/JavaScript-MD5
 * @param  {String} string 加密的字符串
 * @param  {String} key    如果是 HMAC-MD5需要提供key
 * @param  {Boolean} raw    原始块
 * @return {String}        加密后的string或者是raw
 */
qatrix.md5 = function (string, key, raw) {

    /*
     * Add integers, wrapping at 2^32. This uses 16-bit operations internally
     * to work around bugs in some JS interpreters.
     */
    function safeAdd(x, y) {
        var lsw = (x & 0xffff) + (y & 0xffff);
        var msw = (x >> 16) + (y >> 16) + (lsw >> 16);
        return msw << 16 | lsw & 0xffff;
    }

    /*
     * Bitwise rotate a 32-bit number to the left.
     */
    function bitRotateLeft(num, cnt) {
        return num << cnt | num >>> 32 - cnt;
    }

    /*
     * These functions implement the four basic operations the algorithm uses.
     */
    function md5cmn(q, a, b, x, s, t) {
        return safeAdd(bitRotateLeft(safeAdd(safeAdd(a, q), safeAdd(x, t)), s), b);
    }

    function md5ff(a, b, c, d, x, s, t) {
        return md5cmn(b & c | ~b & d, a, b, x, s, t);
    }

    function md5gg(a, b, c, d, x, s, t) {
        return md5cmn(b & d | c & ~d, a, b, x, s, t);
    }

    function md5hh(a, b, c, d, x, s, t) {
        return md5cmn(b ^ c ^ d, a, b, x, s, t);
    }

    function md5ii(a, b, c, d, x, s, t) {
        return md5cmn(c ^ (b | ~d), a, b, x, s, t);
    }

    /*
     * Calculate the MD5 of an array of little-endian words, and a bit length.
     */
    function binlMD5(x, len) {
        /* append padding */
        x[len >> 5] |= 0x80 << len % 32;
        x[(len + 64 >>> 9 << 4) + 14] = len;

        var i;
        var olda;
        var oldb;
        var oldc;
        var oldd;
        var a = 1732584193;
        var b = -271733879;
        var c = -1732584194;
        var d = 271733878;

        for (i = 0; i < x.length; i += 16) {
            olda = a;
            oldb = b;
            oldc = c;
            oldd = d;

            a = md5ff(a, b, c, d, x[i], 7, -680876936);
            d = md5ff(d, a, b, c, x[i + 1], 12, -389564586);
            c = md5ff(c, d, a, b, x[i + 2], 17, 606105819);
            b = md5ff(b, c, d, a, x[i + 3], 22, -1044525330);
            a = md5ff(a, b, c, d, x[i + 4], 7, -176418897);
            d = md5ff(d, a, b, c, x[i + 5], 12, 1200080426);
            c = md5ff(c, d, a, b, x[i + 6], 17, -1473231341);
            b = md5ff(b, c, d, a, x[i + 7], 22, -45705983);
            a = md5ff(a, b, c, d, x[i + 8], 7, 1770035416);
            d = md5ff(d, a, b, c, x[i + 9], 12, -1958414417);
            c = md5ff(c, d, a, b, x[i + 10], 17, -42063);
            b = md5ff(b, c, d, a, x[i + 11], 22, -1990404162);
            a = md5ff(a, b, c, d, x[i + 12], 7, 1804603682);
            d = md5ff(d, a, b, c, x[i + 13], 12, -40341101);
            c = md5ff(c, d, a, b, x[i + 14], 17, -1502002290);
            b = md5ff(b, c, d, a, x[i + 15], 22, 1236535329);

            a = md5gg(a, b, c, d, x[i + 1], 5, -165796510);
            d = md5gg(d, a, b, c, x[i + 6], 9, -1069501632);
            c = md5gg(c, d, a, b, x[i + 11], 14, 643717713);
            b = md5gg(b, c, d, a, x[i], 20, -373897302);
            a = md5gg(a, b, c, d, x[i + 5], 5, -701558691);
            d = md5gg(d, a, b, c, x[i + 10], 9, 38016083);
            c = md5gg(c, d, a, b, x[i + 15], 14, -660478335);
            b = md5gg(b, c, d, a, x[i + 4], 20, -405537848);
            a = md5gg(a, b, c, d, x[i + 9], 5, 568446438);
            d = md5gg(d, a, b, c, x[i + 14], 9, -1019803690);
            c = md5gg(c, d, a, b, x[i + 3], 14, -187363961);
            b = md5gg(b, c, d, a, x[i + 8], 20, 1163531501);
            a = md5gg(a, b, c, d, x[i + 13], 5, -1444681467);
            d = md5gg(d, a, b, c, x[i + 2], 9, -51403784);
            c = md5gg(c, d, a, b, x[i + 7], 14, 1735328473);
            b = md5gg(b, c, d, a, x[i + 12], 20, -1926607734);

            a = md5hh(a, b, c, d, x[i + 5], 4, -378558);
            d = md5hh(d, a, b, c, x[i + 8], 11, -2022574463);
            c = md5hh(c, d, a, b, x[i + 11], 16, 1839030562);
            b = md5hh(b, c, d, a, x[i + 14], 23, -35309556);
            a = md5hh(a, b, c, d, x[i + 1], 4, -1530992060);
            d = md5hh(d, a, b, c, x[i + 4], 11, 1272893353);
            c = md5hh(c, d, a, b, x[i + 7], 16, -155497632);
            b = md5hh(b, c, d, a, x[i + 10], 23, -1094730640);
            a = md5hh(a, b, c, d, x[i + 13], 4, 681279174);
            d = md5hh(d, a, b, c, x[i], 11, -358537222);
            c = md5hh(c, d, a, b, x[i + 3], 16, -722521979);
            b = md5hh(b, c, d, a, x[i + 6], 23, 76029189);
            a = md5hh(a, b, c, d, x[i + 9], 4, -640364487);
            d = md5hh(d, a, b, c, x[i + 12], 11, -421815835);
            c = md5hh(c, d, a, b, x[i + 15], 16, 530742520);
            b = md5hh(b, c, d, a, x[i + 2], 23, -995338651);

            a = md5ii(a, b, c, d, x[i], 6, -198630844);
            d = md5ii(d, a, b, c, x[i + 7], 10, 1126891415);
            c = md5ii(c, d, a, b, x[i + 14], 15, -1416354905);
            b = md5ii(b, c, d, a, x[i + 5], 21, -57434055);
            a = md5ii(a, b, c, d, x[i + 12], 6, 1700485571);
            d = md5ii(d, a, b, c, x[i + 3], 10, -1894986606);
            c = md5ii(c, d, a, b, x[i + 10], 15, -1051523);
            b = md5ii(b, c, d, a, x[i + 1], 21, -2054922799);
            a = md5ii(a, b, c, d, x[i + 8], 6, 1873313359);
            d = md5ii(d, a, b, c, x[i + 15], 10, -30611744);
            c = md5ii(c, d, a, b, x[i + 6], 15, -1560198380);
            b = md5ii(b, c, d, a, x[i + 13], 21, 1309151649);
            a = md5ii(a, b, c, d, x[i + 4], 6, -145523070);
            d = md5ii(d, a, b, c, x[i + 11], 10, -1120210379);
            c = md5ii(c, d, a, b, x[i + 2], 15, 718787259);
            b = md5ii(b, c, d, a, x[i + 9], 21, -343485551);

            a = safeAdd(a, olda);
            b = safeAdd(b, oldb);
            c = safeAdd(c, oldc);
            d = safeAdd(d, oldd);
        }
        return [a, b, c, d];
    }

    /*
     * Convert an array of little-endian words to a string
     */
    function binl2rstr(input) {
        var i;
        var output = '';
        var length32 = input.length * 32;
        for (i = 0; i < length32; i += 8) {
            output += String.fromCharCode(input[i >> 5] >>> i % 32 & 0xff);
        }
        return output;
    }

    /*
     * Convert a raw string to an array of little-endian words
     * Characters >255 have their high-byte silently ignored.
     */
    function rstr2binl(input) {
        var i;
        var output = [];
        output[(input.length >> 2) - 1] = undefined;
        for (i = 0; i < output.length; i += 1) {
            output[i] = 0;
        }
        var length8 = input.length * 8;
        for (i = 0; i < length8; i += 8) {
            output[i >> 5] |= (input.charCodeAt(i / 8) & 0xff) << i % 32;
        }
        return output;
    }

    /*
     * Calculate the MD5 of a raw string
     */
    function rstrMD5(s) {
        return binl2rstr(binlMD5(rstr2binl(s), s.length * 8));
    }

    /*
     * Calculate the HMAC-MD5, of a key and some data (raw strings)
     */
    function rstrHMACMD5(key, data) {
        var i;
        var bkey = rstr2binl(key);
        var ipad = [];
        var opad = [];
        var hash;
        ipad[15] = opad[15] = undefined;
        if (bkey.length > 16) {
            bkey = binlMD5(bkey, key.length * 8);
        }
        for (i = 0; i < 16; i += 1) {
            ipad[i] = bkey[i] ^ 0x36363636;
            opad[i] = bkey[i] ^ 0x5c5c5c5c;
        }
        hash = binlMD5(ipad.concat(rstr2binl(data)), 512 + data.length * 8);
        return binl2rstr(binlMD5(opad.concat(hash), 512 + 128));
    }

    /*
     * Convert a raw string to a hex string
     */
    function rstr2hex(input) {
        var hexTab = '0123456789abcdef';
        var output = '';
        var x;
        var i;
        for (i = 0; i < input.length; i += 1) {
            x = input.charCodeAt(i);
            output += hexTab.charAt(x >>> 4 & 0x0f) + hexTab.charAt(x & 0x0f);
        }
        return output;
    }

    /*
     * Encode a string as utf-8
     */
    function str2rstrUTF8(input) {
        return unescape(encodeURIComponent(input));
    }

    /*
     * Take string arguments and return either raw or hex encoded strings
     */
    function rawMD5(s) {
        return rstrMD5(str2rstrUTF8(s));
    }

    function hexMD5(s) {
        return rstr2hex(rawMD5(s));
    }

    function rawHMACMD5(k, d) {
        return rstrHMACMD5(str2rstrUTF8(k), str2rstrUTF8(d));
    }

    function hexHMACMD5(k, d) {
        return rstr2hex(rawHMACMD5(k, d));
    }

    if (!key) {
        if (!raw) {
            return hexMD5(string);
        }
        return rawMD5(string);
    }
    if (!raw) {
        return hexHMACMD5(key, string);
    }
    return rawHMACMD5(key, string);
};

/**
 * 顶部提示框
 * @param {string} msg 提示信息
 * @param {string} type 提示框类型 success/error
 * @param {number} duration 动画时间
 * @param {number} wait 提示框停留时间
 * @param {function} callback 提示框收起后触发回调函数
 */
qatrix.alert = function (msg, type, duration, wait, callback) {
    var $others = $('.biz-alert-box');
    var $box = $('<div class="biz-alert-box"><span>' + msg + '</span><a></a></div>');
    switch (type) {
        case 'success':
            $box.css({
                'background-color': 'rgba(51, 187, 116, 1)'
            });
            break;
        case 'error':
            $box.css({
                'background-color': 'rgba(239, 59, 44, 1)'
            });
            break;
        default:
            $box.css({
                'background-color': 'rgba(239, 59, 44, 1)'
            });
    }

    // 判断是弹出框还是主页面
    var $container = $('.biz-modal').length > 0 ? $('.biz-modal') : $('.biz-viewport');
    $container.append($box);

    var top = $box.css('top');
    var durationTime = parseInt(duration) || 500;
    var waitTime = parseInt(wait) || 5000;
    // 绑定关闭按钮事件
    $box.find('a').on('click', function () {
        close();
    });
    $box.animate({
        'top': '0'
    }, durationTime, 'swing', function () {
        // 移除已存在的提示框
        $others.remove();
        setTimeout(function () {
            close();
        }, waitTime);
    });
    // 收起提示框
    function close() {
        $box.animate({
            'top': top
        }, durationTime, 'swing', function () {
            $box.remove();
            if (typeof callback == 'function') {
                callback();
            }
        });
    }
};
/**
 * 提示框
 * @param  {[type]} msg [description]
 * @return {[type]}     [description]
 */
qatrix['showToast'] = function (msg, type, duration, wait) {
    var $toast = $('<div class="toast">' + msg + '</div>');
    $('.chat-box').prepend($toast);
    $toast.css({ 'margin-left': -($toast.css('width') / 2) + 'px', 'z-index': '10000' });

    var durationTime = parseInt(duration) || 400;
    var waitTime = parseInt(wait) || 2000;
    $toast.animate({ 'opacity': '1' }, durationTime, function () {
        setTimeout(function () {
            $toast.animate({ 'opacity': '0' }, durationTime, function () {
                $('.toast').remove();
            });
        }, waitTime);
    });
};

/**
 * 生成UUID
 * @return {[type]} [description]
 */
qatrix['uuid'] = function () {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
};

/**
 * 封装axiox的 post 请求
 * @param {Object} params 请求参数
 * @param {Function} successCB 成功回调
 * @param {Function} failCB 失败回调
 */

qatrix['axiosPost'] = function (options) {

    var defaultOptions = {
        method: 'post',
        url: '',
        data: {},
        contentType: 'form-data', //form-urlencodeed
        headers: {
            'content-type': 'multipart/form-data',
			'Authorization':'Bearer ' + qatrix.Cookies['get']('adminapi_token')
        },
        withCredentials: true,
        validateStatus: function validateStatus(status) {
            return status >= 200 && status < 300 || status == 401 || status == 403; //201 created , 401 unauthorization
        }
    };

    var mix_option = Object.assign(defaultOptions, options);
    if (mix_option['contentType'] == 'form-data') {

        var bodyFormData = new FormData();
        for (var key in mix_option['data']) {
            bodyFormData.append(key, mix_option['data'][key]);
        }
        mix_option['data'] = bodyFormData;
    } else if (mix_option['contentType'] == 'form-urlencodeed') {
        mix_option['transformRequest'] = [function (data) {
            var ret = '';
            for (var it in data) {
                ret += encodeURIComponent(it) + '=' + encodeURIComponent(data[it]) + '&';
            }
            return ret;
        }];
        mix_option['headers']['content-type'] = 'application/x-www-form-urlencoded;charset=UTF-8';
    } else if (mix_option['contentType'] == 'json') {
        mix_option['headers']['content-type'] = 'application/json';
    }
    return new Promise(function (resolve, reject) {
        axios(mix_option).then(function (response) {
            if (response.status == 401) {
                //qatrix.alert('查询错误', 'error', 500, 1500);
                var domain = window.location.protocol + '//' + window.location.host;
                // window.location.href = domain + '/login.html?callback_url=' + encodeURIComponent(domain + window.location.pathname + window.location.search);
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            } else if (response.status == 403) {
				alert("操作未授权");
                // window.location.pathname = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            }
            resolve(response.data);
        }).catch(function (error) {
            reject(error);
        });
    });
};

/**
 * 封装axiox的 get 请求，只操作url
 * @param {Object} params 请求参数
 * @param {Function} successCB 成功回调
 * @param {Function} failCB 失败回调
 */
qatrix['axiosGet'] = function (options) {
    var defaultOptions = {
        method: 'get',
        url: '',
        headers: {
            'content-type': 'application/x-www-form-urlencoded',
			'Authorization':'Bearer ' + qatrix.Cookies['get']('adminapi_token')
        },
        //  withCredentials:true,
        // todo test
        withCredentials: false,
        validateStatus: function validateStatus(status) {
            return status >= 200 && status < 300 || status == 401 || status == 403; //201 created , 401 unauthorization
        }
    };
    var mix_option = Object.assign(defaultOptions, options);
    return new Promise(function (resolve, reject) {
        axios(mix_option).then(function (response) {
            if (response.status == 401) {
                //qatrix.alert('查询错误', 'error', 500, 1500);
                var domain = window.location.protocol + '//' + window.location.host;
                // window.location.href = domain + '/login.html?callback_url=' + encodeURIComponent(domain + window.location.pathname + window.location.search);
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            } else if (response.status == 403) {
				alert("操作未授权");
                // window.location.pathname = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            }
            resolve(response.data);
        }).catch(function (error) {
            reject(error);
        });
    });
};
/**
 * 封装axiox的 delete 请求，只操作url
 * @param {Object} params 请求参数
 * @param {Function} successCB 成功回调
 * @param {Function} failCB 失败回调
 */
qatrix['axiosDelete'] = function (options) {
    var defaultOptions = {
        method: 'delete',
        url: '',
        headers: {
			'Authorization':'Bearer ' + qatrix.Cookies['get']('adminapi_token')
		},
        withCredentials: true,
        validateStatus: function validateStatus(status) {
            return status >= 200 && status < 300 || status == 401 || status == 403; //201 created , 401 unauthorization
        }
    };
    var mix_option = Object.assign(defaultOptions, options);
    return new Promise(function (resolve, reject) {
        axios(mix_option).then(function (response) {
            if (response.status == 401) {
                //qatrix.alert('查询错误', 'error', 500, 1500);
                var domain = window.location.protocol + '//' + window.location.host;
                // window.location.href = domain + '/login.html?callback_url=' + encodeURIComponent(domain + window.location.pathname);
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            } else if (response.status == 403) {
				alert("操作未授权");
                // window.location.pathname = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            }
            resolve(response.data);
        }).catch(function (error) {
            reject(error);
        });
    });
};

/**
 * 封装axiox的 put 请求
 * @param {Object} params 请求参数
 * @param {Function} successCB 成功回调
 * @param {Function} failCB 失败回调
 */

qatrix['axiosPut'] = function (options) {

    var defaultOptions = {
        method: 'put',
        url: '',
        data: {},
        contentType: 'form-urlencodeed', //form-urlencodeed
        headers: {
            'content-type': 'multipart/form-data',
			'Authorization':'Bearer ' + qatrix.Cookies['get']('adminapi_token')
        },
        withCredentials: true,
        validateStatus: function validateStatus(status) {
            return status >= 200 && status < 300 || status == 401 || status == 403; //201 created , 401 unauthorization
        }
    };

    var mix_option = Object.assign(defaultOptions, options);
    if (mix_option['contentType'] == 'form-data') {

        var bodyFormData = new FormData();
        for (var key in mix_option['data']) {
            bodyFormData.append(key, mix_option['data'][key]);
        }
        mix_option['data'] = bodyFormData;
    } else if (mix_option['contentType'] == 'form-urlencodeed') {
        mix_option['transformRequest'] = [function (data) {
            var ret = '';

            var _loop = function _loop(it) {
                if (data[it] instanceof Array) {
                    data[it].forEach(function (item, idx) {
                        var key = it + '[' + idx + ']';
                        ret += encodeURIComponent(key) + '=' + encodeURIComponent(JSON.stringify(item)) + '&';
                    });
                } else {
                    ret += encodeURIComponent(it) + '=' + encodeURIComponent(data[it]) + '&';
                }
            };

            for (var it in data) {
                _loop(it);
            }
            return ret;
        }];
        mix_option['headers']['content-type'] = 'application/x-www-form-urlencoded;charset=UTF-8';
    } else if (mix_option['contentType'] == 'json') {
        mix_option['headers']['content-type'] = 'application/json';
    }
    return new Promise(function (resolve, reject) {
        axios(mix_option).then(function (response) {
            if (response.status == 401) {
                //qatrix.alert('查询错误', 'error', 500, 1500);
                var domain = window.location.protocol + '//' + window.location.host;
                // window.location.href = domain + '/login.html?callback_url=' + encodeURIComponent(domain + window.location.pathname + window.location.search);
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            } else if (response.status == 403) {
				alert("操作未授权");
                // window.location.pathname = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
				window.location.href = 'http://corev2.influx.io/index.php?controller=sysauth&action=login';
                return;
            }
            resolve(response.data);
        }).catch(function (error) {
            reject(error);
        });
    });
};

/**
 * 封装请求的url的参数
 * example:
 * '/v1/user/{user_id}/tag' => '/v1/user/1234/tag'
 * '/v1/user/{user_id}/tag/{id}' => '/v1/user/1234/tag/321'
 * @param {Object} params 请求参数 url, id1 [,id2,...]
 * @param {Function} successCB 成功回调
 * @param {Function} failCB 失败回调
 */
qatrix['formatGetUrlQuery'] = function (url, id) {
    if (!url || typeof url != 'string') return false;
    var arg = arguments;
    var _arr = url.match(/\{[^\}]+\}/g);

    for (var key in _arr) {
        url = url.replace(_arr[key], arg[Number(key) + 1]);
    }
    return url;
};

/* 封装url查询参数替换封装，不刷新页面替换url
http://guide.mobmobs.com/reports/orders/index.html => http://guide.mobmobs.com/reports/orders/index.html?start-date=2018-08-20&end-date=2018-08-20
params obj {Object} 查询参数的key、value
*/
qatrix['replaceUrlQuery'] = function (url, obj) {
    var search = location.search;
    search.slice(1).split('&');
};

/*
    cookie
*/
qatrix['Cookies'] = {
    converter: function converter() {},
    extend: function extend() {
        var i = 0;
        var result = {};
        for (; i < arguments.length; i++) {
            var attributes = arguments[i];
            for (var key in attributes) {
                result[key] = attributes[key];
            }
        }
        return result;
    },
    decode: function decode(s) {
        return s.replace(/(%[0-9A-Z]{2})+/g, decodeURIComponent);
    },
    set: function set(key, value, attributes) {
        var self = this;
        if (typeof document === 'undefined') {
            return;
        }

        attributes = self.extend({
            path: '/'
        }, {}, attributes);

        if (typeof attributes.expires === 'number') {
            attributes.expires = new Date(new Date() * 1 + attributes.expires * 864e+5);
        }

        // We're using "expires" because "max-age" is not supported by IE
        attributes.expires = attributes.expires ? attributes.expires.toUTCString() : '';

        try {
            var result = JSON.stringify(value);
            if (/^[\{\[]/.test(result)) {
                value = result;
            }
        } catch (e) {}

        value = self.converter.write ? self.converter.write(value, key) : encodeURIComponent(String(value)).replace(/%(23|24|26|2B|3A|3C|3E|3D|2F|3F|40|5B|5D|5E|60|7B|7D|7C)/g, decodeURIComponent);

        key = encodeURIComponent(String(key)).replace(/%(23|24|26|2B|5E|60|7C)/g, decodeURIComponent).replace(/[\(\)]/g, escape);

        var stringifiedAttributes = '';
        for (var attributeName in attributes) {
            if (!attributes[attributeName]) {
                continue;
            }
            stringifiedAttributes += '; ' + attributeName;
            if (attributes[attributeName] === true) {
                continue;
            }
            stringifiedAttributes += '=' + attributes[attributeName].split(';')[0];
        }

        return document.cookie = key + '=' + value + stringifiedAttributes;
    },
    _get: function _get(key, json) {
        var self = this;
        if (typeof document === 'undefined') {
            return;
        }

        var jar = {};
        var cookies = document.cookie ? document.cookie.split('; ') : [];
        var i = 0;

        for (; i < cookies.length; i++) {
            var parts = cookies[i].split('=');
            var cookie = parts.slice(1).join('=');

            if (!json && cookie.charAt(0) === '"') {
                cookie = cookie.slice(1, -1);
            }
            try {
                var name = self.decode(parts[0]);
                cookie = (self.converter.read || self.converter)(cookie, name) || self.decode(cookie);
                if (json) {
                    try {
                        cookie = JSON.parse(cookie);
                    } catch (e) {}
                }

                jar[name] = cookie;

                if (key === name) {
                    break;
                }
            } catch (e) {}
        }

        return key ? jar[key] : jar;
    },
    get: function get(key) {
        var self = this;
        return self._get(key, false /* read as raw */);
    },
    getJSON: function getJSON(key) {
        var self = this;
        return self._get(key, true /* read as json */);
    },
    remove: function remove(key, attributes) {
        var self = this;
        self.set(key, '', self.extend(attributes, {
            expires: -1
        }));
    }

    /**
     * 封装验证手机号
     *
     */
};qatrix['checkForm'] = {
    phone: function phone(_phone) {
        return (/^1[3-9][0-9]{9}$/.test(_phone)
        );
    }

    /**
      * 弹出框函数：弹出一个窗口，浮动在页面上方

      * 参数说明：
      * url(string):弹出框所包含的网页的url
      * width(string):弹出框的宽度，接受两种类型的参数：auto（默认宽度为920px居中），pixel（具体的宽度值，不含'px'字符）
      * height(string):弹出框的高度，接受两种类型的参数：auto（高度自适应，与浏览器窗口上下沿分别相隔20px），pixel（具体的高度值，不含'px'字符）
      * outsideClose(boolean):点击弹出框以外的区域，是否把弹出框关闭，默认false
      * opacity(number):遮蔽层的透明度，0为全透明，1为不透明，默认为0.8
    */
};function openWindow(url, width, height, outsideClose, opacity) {

    //判断窗口高度是否超过屏幕高度
    if ($(window).height() < height) {
        height = $(window).height() - 20;
    }

    //定义遮蔽层默认透明度
    opacity = opacity == undefined ? 0.8 : opacity;

    //增加遮蔽层（需要用到样式#cover）
    $('body').append($("<div id='cover'></div>"));
    //设置遮蔽层高度（内容区域的高度，而不是可视区域的高度）
    $('#cover').height($(document).height() + 800);
    //添加遮蔽层点击事件
    if (outsideClose) $('#cover').click(function () {
        closeWindow();
    });

    //添加弹出框
    $("body").append($("<div id='examWindow'></div>"));

    //设置弹出框样式
    if (width == '' || width == 'auto') width = '920'; //定义默认宽度
    if (height == '' || height == 'auto') height = $(window).height() - 40; //定义默认高度
    //设置宽高度
    $('#examWindow').css('width', width);
    $('#examWindow').css('height', height);
    //左右居中弹出框
    $('#examWindow').css('margin-left', '-' + width / 2 + 'px');
    //$('#examWindow').css('margin-top','-'+height/2+'px');

    //定位弹出框
    var window_top = $(document).scrollTop() + ($(window).height() - height) / 2;
    $('#examWindow').css('top', window_top + 'px');

    // //添加关闭按钮
    // $('#examWindow').append($("<a id='closeButton'>关闭</a>"));
    // //添加关闭事件
    // $('#closeButton').click(function(){closeWindow()});
    //检测浏览器版本
    var version;
    var ua = navigator.userAgent; //获取用户端信息
    var b = ua.indexOf("MSIE "); //检测特殊字符串"MSIE "的位置
    if (b < 0) version = 0;else version = parseFloat(ua.substring(b + 5, ua.indexOf(";", b))); //截取版本号字符串，并转换为数值
    // //实现关闭按钮在ie6下png透明
    // if(version==6){
    //  DD_belatedPNG.fix('a#closeButton');
    // }

    //内嵌网页
    $('#examWindow').append($("<iframe src='" + url + "' frameborder='0' width='100%' height='100%'></iframe>"));
    //添加阴影层
    $("body").append($("<div class='shadow'></div>"));
    //设置阴影层样式
    var shadowWidth = $('#examWindow').width() + 2;
    var shadowHeight = $('#examWindow').height() + 2;
    $('.shadow').css({ 'width': shadowWidth + 'px', 'height': shadowHeight + 'px', 'margin-left': '-' + shadowWidth / 2 + 'px' });
    //定位阴影层
    var shadow_top = window_top - 1;
    $('.shadow').css('top', shadow_top + 'px');

    //显示遮蔽层
    $('#cover').css('opacity', opacity).fadeIn('normal');
    //显示阴影层
    $('.shadow').fadeIn('normal');
    //显示弹出框
    $('#examWindow').fadeIn('normal');

    //若没有指定高度，则浏览器窗口大小变化时，自动调整弹出框和阴影层的高度
    if (height == '' || height == 'auto') {
        $(window).resize(function () {
            $('#examWindow').css('height', $(window).height() - 40 + 'px');
            $('.shadow').css('height', $('#examWindow').height() + 2 + 'px');
        });
    }

    //限制屏幕滚动
    if ($.browser.msie) {
        document.documentElement.style.overflow = "hidden";
    } else {
        document.body.style.overflow = "hidden";
    }

    return false;
}
//关闭弹出框函数
function closeWindow() {
    //移除遮蔽层，阴影层和弹出框
    $('#examWindow').remove();
    $('.shadow').remove();
    $('#cover').remove();

    //恢复屏幕滚动
    if ($.browser.msie) {
        document.documentElement.style.overflow = "auto";
    } else {
        document.body.style.overflow = "auto";
    }

    return false;
}
