const cookie = require('cookie');
const yargs = require('yargs');
const URL = require('url');
const querystring = require('query-string');
const exec = require("child_process").exec;

const object = {};
const hasOwnProperty = object.hasOwnProperty;
const forOwn = (object, callback) => {
    for (const key in object) {
        if (hasOwnProperty.call(object, key)) {
            callback(key, object[key]);
        }
    }
};

const extend = (destination, source) => {
    if (!source) {
        return destination;
    }
    forOwn(source, (key, value) => {
        destination[key] = value;
    });
    return destination;
};

const forEach = (array, callback) => {
    const length = array.length;
    let index = -1;
    while (++index < length) {
        callback(array[index]);
    }
};

const toString = object.toString;
const isArray = Array.isArray;
const isBuffer = Buffer.isBuffer;
const isObject = (value) => {
    // This is a very simple check, but it’s good enough for what we need.
    return toString.call(value) == '[object Object]';
};
const isString = (value) => {
    return typeof value == 'string' ||
        toString.call(value) == '[object String]';
};
const isNumber = (value) => {
    return typeof value == 'number' ||
        toString.call(value) == '[object Number]';
};
const isFunction = (value) => {
    return typeof value == 'function';
};
const isMap = (value) => {
    return toString.call(value) == '[object Map]';
};
const isSet = (value) => {
    return toString.call(value) == '[object Set]';
};

/*--------------------------------------------------------------------------*/

// https://mathiasbynens.be/notes/javascript-escapes#single
const singleEscapes = {
    '"': '\\"',
    '\'': '\\\'',
    '\\': '\\\\',
    '\b': '\\b',
    '\f': '\\f',
    '\n': '\\n',
    '\r': '\\r',
    '\t': '\\t'
    // `\v` is omitted intentionally, because in IE < 9, '\v' == 'v'.
    // '\v': '\\x0B'
};

const regexSingleEscape = /["'\\\b\f\n\r\t]/;
const regexDigit = /[0-9]/;
const regexWhitelist = /[ !#-&\(-\[\]-_a-~]/;

const jsesc = (argument, options) => {
    const increaseIndentation = () => {
        oldIndent = indent;
        ++options.indentLevel;
        indent = options.indent.repeat(options.indentLevel)
    };
    // Handle options
    const defaults = {
        'escapeEverything': false,
        'minimal': false,
        'isScriptContext': false,
        'quotes': 'single',
        'wrap': false,
        'es6': false,
        'json': false,
        'compact': true,
        'lowercaseHex': false,
        'numbers': 'decimal',
        'indent': '\t',
        'indentLevel': 0,
        '__inline1__': false,
        '__inline2__': false
    };
    const json = options && options.json;
    if (json) {
        defaults.quotes = 'double';
        defaults.wrap = true;
    }
    options = extend(defaults, options);
    if (
        options.quotes != 'single' &&
        options.quotes != 'double' &&
        options.quotes != 'backtick'
    ) {
        options.quotes = 'single';
    }
    const quote = options.quotes == 'double' ?
        '"' :
        (options.quotes == 'backtick' ?
                '`' :
                '\''
        );
    const compact = options.compact;
    const lowercaseHex = options.lowercaseHex;
    let indent = options.indent.repeat(options.indentLevel);
    let oldIndent = '';
    const inline1 = options.__inline1__;
    const inline2 = options.__inline2__;
    const newLine = compact ? '' : '\n';
    let result;
    let isEmpty = true;
    const useBinNumbers = options.numbers == 'binary';
    const useOctNumbers = options.numbers == 'octal';
    const useDecNumbers = options.numbers == 'decimal';
    const useHexNumbers = options.numbers == 'hexadecimal';

    if (json && argument && isFunction(argument.toJSON)) {
        argument = argument.toJSON();
    }

    if (!isString(argument)) {
        if (isMap(argument)) {
            if (argument.size == 0) {
                return 'new Map()';
            }
            if (!compact) {
                options.__inline1__ = true;
                options.__inline2__ = false;
            }
            return 'new Map(' + jsesc(Array.from(argument), options) + ')';
        }
        if (isSet(argument)) {
            if (argument.size == 0) {
                return 'new Set()';
            }
            return 'new Set(' + jsesc(Array.from(argument), options) + ')';
        }
        if (isBuffer(argument)) {
            if (argument.length == 0) {
                return 'Buffer.from([])';
            }
            return 'Buffer.from(' + jsesc(Array.from(argument), options) + ')';
        }
        if (isArray(argument)) {
            result = [];
            options.wrap = true;
            if (inline1) {
                options.__inline1__ = false;
                options.__inline2__ = true;
            }
            if (!inline2) {
                increaseIndentation();
            }
            forEach(argument, (value) => {
                isEmpty = false;
                if (inline2) {
                    options.__inline2__ = false;
                }
                result.push(
                    (compact || inline2 ? '' : indent) +
                    jsesc(value, options)
                );
            });
            if (isEmpty) {
                return '[]';
            }
            if (inline2) {
                return '[' + result.join(', ') + ']';
            }
            return '[' + newLine + result.join(',' + newLine) + newLine +
                (compact ? '' : oldIndent) + ']';
        } else if (isNumber(argument)) {
            if (json) {
                // Some number values (e.g. `Infinity`) cannot be represented in JSON.
                return JSON.stringify(argument);
            }
            if (useDecNumbers) {
                return String(argument);
            }
            if (useHexNumbers) {
                let hexadecimal = argument.toString(16);
                if (!lowercaseHex) {
                    hexadecimal = hexadecimal.toUpperCase();
                }
                return '0x' + hexadecimal;
            }
            if (useBinNumbers) {
                return '0b' + argument.toString(2);
            }
            if (useOctNumbers) {
                return '0o' + argument.toString(8);
            }
        } else if (!isObject(argument)) {
            if (json) {
                // For some values (e.g. `undefined`, `function` objects),
                // `JSON.stringify(value)` returns `undefined` (which isn’t valid
                // JSON) instead of `'null'`.
                return JSON.stringify(argument) || 'null';
            }
            return String(argument);
        } else { // it’s an object
            result = [];
            options.wrap = true;
            increaseIndentation();
            forOwn(argument, (key, value) => {
                isEmpty = false;
                result.push(
                    (compact ? '' : indent) +
                    jsesc(key, options) + ':' +
                    (compact ? '' : ' ') +
                    jsesc(value, options)
                );
            });
            if (isEmpty) {
                return '{}';
            }
            return '{' + newLine + result.join(',' + newLine) + newLine +
                (compact ? '' : oldIndent) + '}';
        }
    }

    const string = argument;
    // Loop over each code unit in the string and escape it
    let index = -1;
    const length = string.length;
    result = '';
    while (++index < length) {
        const character = string.charAt(index);
        if (options.es6) {
            const first = string.charCodeAt(index);
            if ( // check if it’s the start of a surrogate pair
                first >= 0xD800 && first <= 0xDBFF && // high surrogate
                length > index + 1 // there is a next code unit
            ) {
                const second = string.charCodeAt(index + 1);
                if (second >= 0xDC00 && second <= 0xDFFF) { // low surrogate
                    // https://mathiasbynens.be/notes/javascript-encoding#surrogate-formulae
                    const codePoint = (first - 0xD800) * 0x400 + second - 0xDC00 + 0x10000;
                    let hexadecimal = codePoint.toString(16);
                    if (!lowercaseHex) {
                        hexadecimal = hexadecimal.toUpperCase();
                    }
                    result += '\\u{' + hexadecimal + '}';
                    ++index;
                    continue;
                }
            }
        }
        if (!options.escapeEverything) {
            if (regexWhitelist.test(character)) {
                // It’s a printable ASCII character that is not `"`, `'` or `\`,
                // so don’t escape it.
                result += character;
                continue;
            }
            if (character == '"') {
                result += quote == character ? '\\"' : character;
                continue;
            }
            if (character == '`') {
                result += quote == character ? '\\`' : character;
                continue;
            }
            if (character == '\'') {
                result += quote == character ? '\\\'' : character;
                continue;
            }
        }
        if (
            character == '\0' &&
            !json &&
            !regexDigit.test(string.charAt(index + 1))
        ) {
            result += '\\0';
            continue;
        }
        if (regexSingleEscape.test(character)) {
            // no need for a `hasOwnProperty` check here
            result += singleEscapes[character];
            continue;
        }
        const charCode = character.charCodeAt(0);
        if (options.minimal && charCode != 0x2028 && charCode != 0x2029) {
            result += character;
            continue;
        }
        let hexadecimal = charCode.toString(16);
        if (!lowercaseHex) {
            hexadecimal = hexadecimal.toUpperCase();
        }
        const longhand = hexadecimal.length > 2 || json;
        const escaped = '\\' + (longhand ? 'u' : 'x') +
            ('0000' + hexadecimal).slice(longhand ? -4 : -2);
        result += escaped;
        continue;
    }
    if (options.wrap) {
        result = quote + result + quote;
    }
    if (quote == '`') {
        result = result.replace(/\$\{/g, '\\\$\{');
    }
    if (options.isScriptContext) {
        // https://mathiasbynens.be/notes/etago
        return result
            .replace(/<\/(script|style)/gi, '<\\/$1')
            .replace(/<!--/g, json ? '\\u003C!--' : '\\x3C!--');
    }
    return result;
};

const parseCurlCommand = curlCommand => {
    curlCommand = curlCommand.replace(/\\\r|\\\n/g, '');

    curlCommand = curlCommand.replace(/\s+/g, ' ');

    curlCommand = curlCommand.replace(/ -XPOST/, ' -X POST');
    curlCommand = curlCommand.replace(/ -XGET/, ' -X GET');
    curlCommand = curlCommand.replace(/ -XPUT/, ' -X PUT');
    curlCommand = curlCommand.replace(/ -XPATCH/, ' -X PATCH');
    curlCommand = curlCommand.replace(/ -XDELETE/, ' -X DELETE');
    curlCommand = curlCommand.replace(/ -Xnull/, ' ');
    curlCommand = curlCommand.trim();

    const parsedArguments = yargs
        .boolean(['I', 'head', 'compressed', 'L', 'k', 'silent', 's'])
        .alias('H', 'header')
        .alias('A', 'user-agent')
        .parse(curlCommand);

    let proxy = curlCommand.indexOf('--proxy ') > -1 ?
        curlCommand.match(/--proxy '.*?'/g)[0].replace('--proxy ', '').replaceAll("'", "") :
        curlCommand.indexOf('-x ') > -1 ?
            curlCommand.match(/-x '.*?'/g)[0].replace('-x ', '').replaceAll("'", "") : null;

    let cookieString;
    let cookies;
    let url = parsedArguments._[1].replaceAll("'", "");

    if (!url) {
        for (const argName in parsedArguments) {
            if (parsedArguments.hasOwnProperty(argName) && typeof parsedArguments[argName] === 'string') {
                if (parsedArguments[argName].indexOf('http') === 0 || parsedArguments[argName].indexOf('www.') === 0) {
                    url = parsedArguments[argName];
                }
            }
        }
    }

    let headers;

    if (parsedArguments.header) {
        if (!headers) {
            headers = {};
        }
        if (!Array.isArray(parsedArguments.header)) {
            parsedArguments.header = [parsedArguments.header];
        }
        parsedArguments.header.forEach(header => {
            if (header.indexOf('Cookie') !== -1) {
                cookieString = header;
            } else {
                const components = header.split(/:(.*)/);
                if (components[1]) {
                    headers[components[0]] = components[1].trim();
                }
            }
        });
    }

    if (parsedArguments['user-agent']) {
        if (!headers) {
            headers = {};
        }
        headers['User-Agent'] = parsedArguments['user-agent'];
    }

    if (parsedArguments.b) {
        cookieString = parsedArguments.b;
    }
    if (parsedArguments.cookie) {
        cookieString = parsedArguments.cookie;
    }
    let multipartUploads;
    if (parsedArguments.F) {
        multipartUploads = {};
        if (!Array.isArray(parsedArguments.F)) {
            parsedArguments.F = [parsedArguments.F];
        }
        parsedArguments.F.forEach(multipartArgument => {
            const splitArguments = multipartArgument.split('=', 2);
            const key = splitArguments[0];
            multipartUploads[key] = splitArguments[1];
        });
    }
    if (cookieString) {
        const cookieParseOptions = {
            decode: function (s) {
                return s;
            }
        };
        cookies = cookie.parse(cookieString.replace(/^Cookie: /gi, ''), cookieParseOptions);
    }
    let method;
    if (parsedArguments.X === 'POST') {
        method = 'post';
    } else if (parsedArguments.X === 'PUT' ||
        parsedArguments.T) {
        method = 'put';
    } else if (parsedArguments.X === 'PATCH') {
        method = 'patch';
    } else if (parsedArguments.X === 'DELETE') {
        method = 'delete';
    } else if (parsedArguments.X === 'OPTIONS') {
        method = 'options';
    } else if ((parsedArguments.d ||
        parsedArguments.data ||
        parsedArguments['data-ascii'] ||
        parsedArguments['data-binary'] ||
        parsedArguments['data-raw'] ||
        parsedArguments.F ||
        parsedArguments.form) && !((parsedArguments.G || parsedArguments.get))) {
        method = 'post';
    } else if (parsedArguments.I ||
        parsedArguments.head) {
        method = 'head';
    } else {
        method = 'get';
    }

    const compressed = !!parsedArguments.compressed;
    const urlObject = URL.parse(url);

    if (parsedArguments.G || parsedArguments.get) {
        urlObject.query = urlObject.query ? urlObject.query : '';
        const option = 'd' in parsedArguments ? 'd' : 'data' in parsedArguments ? 'data' : null;
        if (option) {
            let urlQueryString = '';

            if (url.indexOf('?') < 0) {
                url += '?';
            } else {
                urlQueryString += '&';
            }

            if (typeof (parsedArguments[option]) === 'object') {
                urlQueryString += parsedArguments[option].join('&');
            } else {
                urlQueryString += parsedArguments[option];
            }
            urlObject.query += urlQueryString;
            url += urlQueryString;
            delete parsedArguments[option];
        }
    }
    if (urlObject.query && urlObject.query.endsWith('&')) {
        urlObject.query = urlObject.query.slice(0, -1);
    }
    const query = querystring.parse(urlObject.query, {sort: false});
    for (const param in query) {
        if (query[param] === null) {
            query[param] = '';
        }
    }

    urlObject.search = null;
    const request = {
        url: url,
        urlWithoutQuery: URL.format(urlObject)
    };
    if (compressed) {
        request.compressed = true;
    }

    if (Object.keys(query).length > 0) {
        request.query = query;
    }
    if (headers) {
        request.headers = headers;
    }
    request.method = method;

    if (cookies) {
        request.cookies = cookies;
        request.cookieString = cookieString.replace('Cookie: ', '');
    }
    if (multipartUploads) {
        request.multipartUploads = multipartUploads;
    }
    if (parsedArguments.data) {
        request.data = parsedArguments.data;
    } else if (parsedArguments['data-binary']) {
        request.data = parsedArguments['data-binary'];
        request.isDataBinary = true;
    } else if (parsedArguments.d) {
        request.data = parsedArguments.d;
    } else if (parsedArguments['data-ascii']) {
        request.data = parsedArguments['data-ascii'];
    } else if (parsedArguments['data-raw']) {
        request.data = parsedArguments['data-raw'];
        request.isDataRaw = true;
    }

    if (parsedArguments.u) {
        request.auth = parsedArguments.u;
    }
    if (parsedArguments.user) {
        request.auth = parsedArguments.user;
    }
    if (Array.isArray(request.data)) {
        request.dataArray = request.data;
        request.data = request.data.join('&');
    }

    if (parsedArguments.k || parsedArguments.insecure) {
        request.insecure = true;
    }
    if (proxy) {
        request.proxy = proxy;
    }
    return request;
};

const capitalizeUpper = str => {
    return str.toLowerCase().replace(/( |^)[a-z]/g, (L) => L.toUpperCase());
};

const jsonIndent = (json, indent) => {
    indent = indent ? indent : '\t';
    let jsonStr = JSON.stringify(json);
    let result = jsonStr.replaceAll(',"', `,\n${indent}"`).replace('{"', `{\n${indent}"`).replace('"}', `",\n${indent.substr(0, indent.length - 1)}}`);
    let longest = 0;
    let keySet = [];
    Object.keys(json).forEach(key => {
        keySet.push(key);
        longest = key.length > longest ? key.length : longest;
    });
    for (key of keySet) {
        result = padSpace(result, `"${key}"`, longest - key.length + 1);
    }
    return result
};

const str2json = str => {
    let result = {};
    str.split("&").forEach(v => {
        let kv = v.split("=");
        result[kv[0]] = kv[1];
    });
    return JSON.stringify(result);
};

const padSpace = (str, key, num) => {
    for (var i = 0; i < num; i++) {
        str = str.replace(`${key}:`, `${key}:` + " ");
    }
    return str;
};

const serializeCookies = cookieDict => {
    let cookieString = ''
    let i = 0
    const cookieCount = Object.keys(cookieDict).length
    for (const cookieName in cookieDict) {
        const cookieValue = cookieDict[cookieName]
        cookieString += cookieName + '=' + cookieValue
        if (i < cookieCount - 1) {
            cookieString += '; '
        }
        i++
    }
    return cookieString
}

function toGoRequests(curlCommand, blank = "all") {
    const request = parseCurlCommand(curlCommand);
    let code = "";
    if (blank === "all") {
        code += 'package main\n\n';
        code += 'import (\n\t"github.com/Esbiya/requests"\n\t"log"\n)\n\n';
        code += 'func main() {\n';
    }
    code += `\turl := "${request.url}"\n`;
    let requestLineBody = '\tresp := requests.' + capitalizeUpper(request.method) + "(url";
    if (request.headers) {
        code += `\theaders := requests.Headers${jsonIndent(request.headers, '\t\t')}\n`;
        requestLineBody += ", headers";
    }
    if (request.cookies) {
        code += `\tcookies := requests.SimpleCookie${jsonIndent(request.cookies, '\t\t')}\n`;
        requestLineBody += ", cookies";
    }
    if (request.auth) {
        const splitAuth = request.auth.split(':');
        const user = splitAuth[0] || '';
        const password = splitAuth[1] || '';
        code += `\tauth := requests.Auth${jsonIndent({username: user, password: password}, '\t\t')}\n`;
        requestLineBody += ", auth";
    }
    if (request.data === true) {
        request.data = '';
    }
    if (request.data) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString();
        }
        var data;

        try {
            data = JSON.parse(request.data);
            code += `\tdata := requests.Form${jsonIndent(data, '\t\t')}\n`;
        } catch (e) {
            if (request.data.indexOf("&") > -1) {
                data = JSON.parse(str2json(request.data));
                code += `\tdata := requests.Payload${jsonIndent(data, '\t\t')}\n`;
            } else {
                code += `\tdata := []byte("${request.data}")\n`
            }
        }
        requestLineBody += ", data";
    }
    let args1 = ", requests.Arguments{\n";
    let argsLongest;
    if (request.proxy) {
        args1 += `\t\tProxy: "${request.proxy}",\n`;
        argsLongest = 5;
    }
    if (request.insecure) {
        args1 += "\t\tSkipVerifyTLS: true,\n";
        args1 = padSpace(args, "Proxy", 13 - argsLongest);
    }
    args1.length > 22 && (requestLineBody += args1);
    requestLineBody += '})\n';
    code += requestLineBody;
    code += '\tif resp.Error() != nil {\n\t\tlog.Fatal(resp.Error())\n\t}\n';
    code += '\tlog.Println(resp.Text)\n';
    if (blank === "all") code += '}';
    return code + '\n';
}

const toJSONString = function (curlCommand) {
    function repr (value, isKey) {
        return isKey ? "'" + jsesc(value, { quotes: 'single' }) + "'" : value
    }

    function getQueries (request) {
        const queries = {}
        for (const paramName in request.query) {
            const rawValue = request.query[paramName]
            let paramValue
            if (Array.isArray(rawValue)) {
                paramValue = rawValue.map(repr)
            } else {
                paramValue = repr(rawValue)
            }
            queries[repr(paramName)] = paramValue
        }

        return queries
    }

    function getDataString (request) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString()
        }

        const parsedQueryString = querystring.parse(request.data, { sort: false })
        const keyCount = Object.keys(parsedQueryString).length
        const singleKeyOnly = keyCount === 1 && !parsedQueryString[Object.keys(parsedQueryString)[0]]
        const singularData = request.isDataBinary || singleKeyOnly
        if (singularData) {
            const data = {}
            data[repr(request.data)] = ''
            return { data: data }
        } else {
            return getMultipleDataString(request, parsedQueryString)
        }
    }

    function getMultipleDataString (request, parsedQueryString) {
        const data = {}

        for (const key in parsedQueryString) {
            const value = parsedQueryString[key]
            if (Array.isArray(value)) {
                data[repr(key)] = value
            } else {
                data[repr(key)] = repr(value)
            }
        }

        return { data: data }
    }

    function getFilesString (request) {
        const data = {}

        data.files = {}
        data.data = {}

        for (const multipartKey in request.multipartUploads) {
            const multipartValue = request.multipartUploads[multipartKey]
            if (multipartValue.startsWith('@')) {
                const fileName = multipartValue.slice(1)
                data.files[repr(multipartKey)] = repr(fileName)
            } else {
                data.data[repr(multipartKey)] = repr(multipartValue)
            }
        }

        if (Object.keys(data.files).length === 0) {
            delete data.files
        }

        if (Object.keys(data.data).length === 0) {
            delete data.data
        }

        return data
    }

    const request = parseCurlCommand(curlCommand);

    const requestJson = {};

    if (!request.url.match(/https?:/)) {
        request.url = 'http://' + request.url;
    }
    if (!request.urlWithoutQuery.match(/https?:/)) {
        request.urlWithoutQuery = 'http://' + request.urlWithoutQuery;
    }

    requestJson.url = request.urlWithoutQuery.replace(/\/$/, '');
    requestJson.raw_url = request.url;
    requestJson.method = request.method;

    if (request.cookies) {
        const cookies = {};
        for (const cookieName in request.cookies) {
            cookies[repr(cookieName)] = repr(request.cookies[cookieName]);
        }

        requestJson.cookies = cookies;
    }

    if (request.headers) {
        const headers = {};
        for (const headerName in request.headers) {
            headers[repr(headerName)] = repr(request.headers[headerName]);
        }

        requestJson.headers = headers;
    }

    if (request.query) {
        requestJson.queries = getQueries(request);
    }

    if (typeof request.data === 'string' || typeof request.data === 'number') {
        Object.assign(requestJson, getDataString(request));
    } else if (request.multipartUploads) {
        Object.assign(requestJson, getFilesString(request));
    }

    if (request.insecure) {
        requestJson.insecure = false;
    }

    if (request.proxy) {
        requestJson.proxy = request.proxy;
    }

    if (request.auth) {
        const splitAuth = request.auth.split(':');
        const user = splitAuth[0] || '';
        const password = splitAuth[1] || '';

        requestJson.auth = {
            user: repr(user),
            password: repr(password)
        }
    }

    return JSON.stringify(Object.keys(requestJson).length ? requestJson : '{}', null, 4) + '\n';
}

const toPython = function (curlCommand, async = false) {
    function reprWithVariable(value, hasEnvironmentVariable){
        if (!value) {
            return "''";
        }

        if (!hasEnvironmentVariable){
            return "'" + jsesc(value, { quotes: 'single' }) + "'";
        }

        return 'f"' + jsesc(value, { quotes: 'double'}) + '"';
    }

    function repr (value) {
        // In context of url parameters, don't accept nulls and such.
        return reprWithVariable(value, false);
    }

    function getQueryDict (request) {
        let queryDict = 'params = (\n'
        for (const paramName in request.query) {
            const rawValue = request.query[paramName]
            let paramValue
            if (Array.isArray(rawValue)) {
                paramValue = '[' + rawValue.map(repr).join(', ') + ']'
            } else {
                paramValue = repr(rawValue)
            }
            queryDict += '    (' + repr(paramName) + ', ' + paramValue + '),\n'
        }
        queryDict += ')\n'
        return queryDict
    }

    function getDataString (request) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString()
        }
        if (!request.isDataRaw && request.data.startsWith('@')) {
            const filePath = request.data.slice(1)
            if (request.isDataBinary) {
                return 'data = open(\'' + filePath + '\', \'rb\').read()'
            } else {
                return 'data = open(\'' + filePath + '\')'
            }
        }

        const parsedQueryString = querystring.parse(request.data, { sort: false })
        const keyCount = Object.keys(parsedQueryString).length
        const singleKeyOnly = keyCount === 1 && !parsedQueryString[Object.keys(parsedQueryString)[0]]
        const singularData = request.isDataBinary || singleKeyOnly
        if (singularData) {
            return 'data = ' + repr(request.data) + '\n'
        } else {
            return getMultipleDataString(request, parsedQueryString)
        }
    }

    function getMultipleDataString (request, parsedQueryString) {
        let repeatedKey = false
        for (const key in parsedQueryString) {
            const value = parsedQueryString[key]
            if (Array.isArray(value)) {
                repeatedKey = true
            }
        }

        let dataString
        if (repeatedKey) {
            dataString = 'data = [\n'
            for (const key in parsedQueryString) {
                const value = parsedQueryString[key]
                if (Array.isArray(value)) {
                    for (let i = 0; i < value.length; i++) {
                        dataString += '  (' + repr(key) + ', ' + repr(value[i]) + '),\n'
                    }
                } else {
                    dataString += '  (' + repr(key) + ', ' + repr(value) + '),\n'
                }
            }
            dataString += ']\n'
        } else {
            dataString = 'data = {\n'
            const elementCount = Object.keys(parsedQueryString).length
            let i = 0
            for (const key in parsedQueryString) {
                const value = parsedQueryString[key]
                dataString += '  ' + repr(key) + ': ' + repr(value)
                if (i === elementCount - 1) {
                    dataString += '\n'
                } else {
                    dataString += ',\n'
                }
                ++i
            }
            dataString += '}\n'
        }

        return dataString
    }

    function getFilesString (request) {
        // http://docs.python-requests.org/en/master/user/quickstart/#post-a-multipart-encoded-file
        let filesString = 'files = {\n'
        for (const multipartKey in request.multipartUploads) {
            const multipartValue = request.multipartUploads[multipartKey]
            if (multipartValue.startsWith('@')) {
                const fileName = multipartValue.slice(1)
                filesString += '    ' + repr(multipartKey) + ': (' + repr(fileName) + ', open(' + repr(fileName) + ", 'rb')),\n"
            } else {
                filesString += '    ' + repr(multipartKey) + ': (None, ' + repr(multipartValue) + '),\n'
            }
        }
        filesString += '}\n'

        return filesString
    }

    function detectEnvVar(inputString){
        const IN_ENV_VAR = 0, IN_STRING = 1;

        // We only care for the unique element
        let detectedVariables = new Set();
        let currState = IN_STRING;
        let envVarStartIndex = -1;

        const whiteSpaceSet = new Set();
        whiteSpaceSet.add(' ');
        whiteSpaceSet.add('\n');
        whiteSpaceSet.add('\t');

        let modifiedString = [];
        for(const idx in inputString){
            const currIdx = +idx;
            const currChar = inputString[currIdx];
            if(currState === IN_ENV_VAR && whiteSpaceSet.has(currChar)){
                const newVariable = inputString.substring(envVarStartIndex, currIdx);

                if (newVariable !== "") {
                    detectedVariables.add(newVariable);

                    // Change $ -> {
                    // Add } after the last variable name
                    modifiedString.push("{"+newVariable+"}" + currChar);
                }
                else {
                    modifiedString.push("$" + currChar);
                }
                currState = IN_STRING;
                envVarStartIndex = -1;
                continue
            }

            if(currState === IN_ENV_VAR){
                // Skip until we actually have the new variable
                continue;
            }

            // currState === IN_STRING
            if(currChar === '$'){
                currState = IN_ENV_VAR;
                envVarStartIndex = currIdx + 1;
            } else {
                modifiedString.push(currChar);
            }
        }

        if(currState === IN_ENV_VAR){
            const newVariable = inputString.substring(envVarStartIndex, inputString.length);

            if (newVariable !== "") {
                detectedVariables.add(newVariable);
                modifiedString.push("{"+newVariable+"}");
            }
            else {
                modifiedString.push("$");
            }
        }

        return [detectedVariables, modifiedString.join('')];
    }

    const request = parseCurlCommand(curlCommand)

    let osVariables = new Set();

    let cookieDict;
    if (request.cookies) {
        cookieDict = 'cookies = {\n';
        for (const cookieName in request.cookies) {
            const [detectedVars, modifiedString] = detectEnvVar(request.cookies[cookieName]);

            const hasEnvironmentVariable = detectedVars.size > 0;

            for(const newVar of detectedVars){
                osVariables.add(newVar);
            }

            cookieDict += '    ' + repr(cookieName) + ': ' + reprWithVariable(modifiedString, hasEnvironmentVariable) + ',\n';
        }
        cookieDict += '}\n';
    }
    let headerDict;
    if (request.headers) {
        headerDict = 'headers = {\n';
        for (const headerName in request.headers) {
            const [detectedVars, modifiedString] = detectEnvVar(request.headers[headerName]);

            const hasVariable = detectedVars.size > 0;

            for(const newVar of detectedVars){
                osVariables.add(newVar);
            }

            headerDict += '    ' + repr(headerName) + ': ' + reprWithVariable(modifiedString, hasVariable) + ',\n';
        }
        headerDict += '}\n';
    }
    let proxyDict;
    if (request.proxy) {
        proxyDict = async ? `proxy = '${request.proxy}'` : `proxies = {\n\t'http': '${request.proxy}',\n\t'https': '${request.proxy}',\n}\n`;
    }

    let queryDict;
    if (request.query) {
        queryDict = getQueryDict(request);
    }

    let dataString;
    let filesString;
    if (typeof request.data === 'string' || typeof request.data === 'number') {
        dataString = getDataString(request);
    } else if (request.multipartUploads) {
        filesString = getFilesString(request);
    }
    if (!request.url.match(/https?:/)) {
        request.url = 'http://' + request.url;
    }
    if (!request.urlWithoutQuery.match(/https?:/)) {
        request.urlWithoutQuery = 'http://' + request.urlWithoutQuery;
    }
    var requestFirst;
    async ? requestFirst = "async with aiohttp.ClientSession() as session:\n    async with session." : requestFirst = "response = requests.";
    let requestLineWithUrlParams = requestFirst + request.method + '(\'' + request.urlWithoutQuery + '\'';
    let requestLineWithOriginalUrl = requestFirst + request.method + '(\'' + request.url + '\'';

    let requestLineBody = '';
    if (request.headers) {
        requestLineBody += ', headers=headers';
    }
    if (request.query) {
        requestLineBody += ', params=params';
    }
    if (request.cookies) {
        requestLineBody += ', cookies=cookies';
    }
    if (request.proxy) {
        async ? requestLineBody += ', proxy=proxy' : requestLineBody += ', proxies=proxies';
    }
    if (typeof request.data === 'string') {
        requestLineBody += ', data=data';
    } else if (request.multipartUploads) {
        requestLineBody += ', files=files';
    }
    if (request.insecure) {
        async ? requestLineBody += ', ssl=False' : requestLineBody += ', verify=False';
    }
    if (request.auth) {
        const splitAuth = request.auth.split(':');
        const user = splitAuth[0] || '';
        const password = splitAuth[1] || '';
        requestLineBody += ', auth=(' + repr(user) + ', ' + repr(password) + ')';
    }
    requestLineBody += ')';

    if (async) requestLineBody += " as response:\n\tresp = await response.text()";

    requestLineWithOriginalUrl += requestLineBody.replace(', params=params', '');
    requestLineWithUrlParams += requestLineBody;

    let pythonCode = '';

    if (osVariables.size > 0) {
        pythonCode += 'import os\n';
    }

    async ? pythonCode += 'import aiohttp\n\n' : pythonCode += 'import requests\n\n';

    if (osVariables.size > 0) {
        for(const osVar of osVariables){
            const line = `${osVar} = os.getenv('${osVar}')\n`;
            pythonCode += line;
        }

        pythonCode += '\n';
    }

    if (cookieDict) {
        pythonCode += cookieDict + '\n';
    }
    if (headerDict) {
        pythonCode += headerDict + '\n';
    }
    if (proxyDict) {
        pythonCode += proxyDict + '\n';
    }
    if (queryDict) {
        pythonCode += queryDict + '\n';
    }
    if (dataString) {
        pythonCode += dataString + '\n';
    } else if (filesString) {
        pythonCode += filesString + '\n';
    }
    pythonCode += requestLineWithUrlParams;

    return pythonCode + '\n';
}

const toGolang = curlCommand => {
    const request = parseCurlCommand(curlCommand);
    let goCode = 'package main\n\n';
    goCode += 'import (\n\t"fmt"\n\t"io/ioutil"\n\t"log"\n\t"net/http"\n)\n\n';
    goCode += 'func main() {\n';
    goCode += '\tclient := &http.Client{}\n';
    if (request.proxy) {
        goCode += `\tproxy, err := url.Parse("${request.proxy}")\n`;
        goCode += '\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n';
        goCode += '\tclient.Transport = &http.Transport{Proxy: http.ProxyUrl(proxy)}\n';
    }
    if (request.data === true) {
        request.data = '';
    }
    if (request.data) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString();
        }
        if (request.data.indexOf("'") > -1) {
            request.data = jsesc(request.data);
        }
        goCode = goCode.replace('\n)', '\n\t"strings"\n)');
        goCode += '\tvar data = strings.NewReader(`' + request.data + '`)\n';
        goCode += '\treq, err := http.NewRequest("' + request.method.toUpperCase() + '", "' + request.url + '", data)\n';
    } else {
        goCode += '\treq, err := http.NewRequest("' + request.method.toUpperCase() + '", "' + request.url + '", nil)\n';
    }
    goCode += '\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n';
    if (request.headers || request.cookies) {
        for (const headerName in request.headers) {
            goCode += '\treq.Header.Set("' + headerName + '", "' + request.headers[headerName] + '")\n';
        }
        if (request.cookies) {
            const cookieString = serializeCookies(request.cookies);
            goCode += '\treq.Header.Set("Cookie", "' + cookieString + '")\n';
        }
    }

    if (request.auth) {
        const splitAuth = request.auth.split(':');
        const user = splitAuth[0] || '';
        const password = splitAuth[1] || '';
        goCode += '\treq.SetBasicAuth("' + user + '", "' + password + '")\n';
    }
    goCode += '\tresp, err := client.Do(req)\n';
    goCode += '\tif err != nil {\n';
    goCode += '\t\tlog.Fatal(err)\n';
    goCode += '\t}\n';
    goCode += '\tbodyText, err := ioutil.ReadAll(resp.Body)\n';
    goCode += '\tif err != nil {\n';
    goCode += '\t\tlog.Fatal(err)\n';
    goCode += '\t}\n';
    goCode += '\tfmt.Printf("%s\\n", bodyText)\n';
    goCode += '}';

    return goCode + '\n';
}

const toNodeRequest = curlCommand => {
    const request = parseCurlCommand(curlCommand);
    let nodeRequestCode = 'var request = require(\'request\');\n\n';
    if (request.headers || request.cookies) {
        nodeRequestCode += 'var headers = {\n';
        const headerCount = Object.keys(request.headers).length;
        let i = 0;
        for (const headerName in request.headers) {
            nodeRequestCode += '    \'' + headerName + '\': \'' + request.headers[headerName] + '\'';
            if (i < headerCount - 1 || request.cookies) {
                nodeRequestCode += ',\n';
            } else {
                nodeRequestCode += '\n';
            }
            i++;
        }
        if (request.cookies) {
            const cookieString = serializeCookies(request.cookies)
            nodeRequestCode += '    \'Cookie\': \'' + cookieString + '\'\n';
        }
        nodeRequestCode += '};\n\n';
    }

    if (request.data === true) {
        request.data = '';
    }
    if (request.data) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString();
        }
        if (request.data.indexOf("'") > -1) {
            request.data = jsesc(request.data);
        }
        nodeRequestCode += 'var dataString = \'' + request.data + '\';\n\n';
    }

    nodeRequestCode += 'var options = {\n';
    nodeRequestCode += '    url: \'' + request.url + '\'';
    if (request.method !== 'get') {
        nodeRequestCode += ',\n    method: \'' + request.method.toUpperCase() + '\'';
    }

    if (request.headers || request.cookies) {
        nodeRequestCode += ',\n';
        nodeRequestCode += '    headers: headers';
    }
    if (request.data) {
        nodeRequestCode += ',\n    body: dataString';
    }

    if (request.auth) {
        nodeRequestCode += ',\n';
        const splitAuth = request.auth.split(':');
        const user = splitAuth[0] || '';
        const password = splitAuth[1] || '';
        nodeRequestCode += '    auth: {\n';
        nodeRequestCode += "        'user': '" + user + "',\n";
        nodeRequestCode += "        'pass': '" + password + "'\n";
        nodeRequestCode += '    }\n';
    } else {
        nodeRequestCode += '\n';
    }
    nodeRequestCode += '};\n\n';

    nodeRequestCode += 'function callback(error, response, body) {\n';
    nodeRequestCode += '    if (!error && response.statusCode == 200) {\n';
    nodeRequestCode += '        console.log(body);\n';
    nodeRequestCode += '    }\n';
    nodeRequestCode += '}\n\n';
    nodeRequestCode += 'request(options, callback);';

    return nodeRequestCode + '\n';
}

const toJsFetch = curlCommand => {
    const request = parseCurlCommand(curlCommand);

    let jsFetchCode = '';

    if (request.data === true) {
        request.data = '';
    }
    if (request.data) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString();
        }
        if (request.data.indexOf("'") > -1) {
            request.data = jsesc(request.data);
        }

        try {
            JSON.parse(request.data);

            if (!request.headers) {
                request.headers = {};
            }

            if (!request.headers['Content-Type']) {
                request.headers['Content-Type'] = 'application/json; charset=UTF-8';
            }

            request.data = 'JSON.stringify(' + request.data + ')';
        } catch {
            request.data = '\'' + request.data + '\'';
        }
    }

    jsFetchCode += 'fetch(\'' + request.url + '\'';

    if (request.method !== 'get' || request.headers || request.cookies || request.auth || request.body) {
        jsFetchCode += ', {\n';

        if (request.method !== 'get') {
            jsFetchCode += '    method: \'' + request.method.toUpperCase() + '\'';
        }

        if (request.headers || request.cookies || request.auth) {
            if (request.method !== 'get') {
                jsFetchCode += ',\n';
            }
            jsFetchCode += '    headers: {\n';
            const headerCount = Object.keys(request.headers || {}).length;
            let i = 0;
            for (const headerName in request.headers) {
                jsFetchCode += '        \'' + headerName + '\': \'' + request.headers[headerName] + '\'';
                if (i < headerCount - 1 || request.cookies || request.auth) {
                    jsFetchCode += ',\n';
                }
                i++;
            }
            if (request.auth) {
                const splitAuth = request.auth.split(':');
                const user = splitAuth[0] || '';
                const password = splitAuth[1] || '';
                jsFetchCode += '        \'Authorization\': \'Basic \' + btoa(\'' + user + ':' + password + '\')';
            }
            if (request.cookies) {
                const cookieString = serializeCookies(request.cookies);
                jsFetchCode += '        \'Cookie\': \'' + cookieString + '\'';
            }

            jsFetchCode += '\n    }';
        }

        if (request.data) {
            jsFetchCode += ',\n    body: ' + request.data;
        }

        jsFetchCode += '\n}';
    }

    jsFetchCode += ');';

    return jsFetchCode + '\n';
}

const toNodeFetch = curlCommand => {
    let nodeFetchCode = 'var fetch = require(\'node-fetch\');\n\n';
    nodeFetchCode += toJsFetch(curlCommand);

    return nodeFetchCode;
}

const toPhp = curlCommand => {
    const quote = str => jsesc(str, { quotes: 'single' });

    const request = parseCurlCommand(curlCommand);

    let headerString = false;
    if (request.headers) {
        headerString = '$headers = array(\n';
        let i = 0;
        const headerCount = Object.keys(request.headers).length;
        for (const headerName in request.headers) {
            headerString += "    '" + headerName + "' => '" + quote(request.headers[headerName]) + "'";
            if (i < headerCount - 1) {
                headerString += ',\n';
            }
            i++;
        }
        if (request.cookies) {
            const cookieString = quote(serializeCookies(request.cookies));
            headerString += ",\n    'Cookie' => '" + cookieString + "'";
        }
        headerString += '\n);';
    } else {
        headerString = '$headers = array();';
    }

    let optionsString = false;
    if (request.auth) {
        const splitAuth = request.auth.split(':').map(quote);
        const user = splitAuth[0] || '';
        const password = splitAuth[1] || '';
        optionsString = "$options = array('auth' => array('" + user + "', '" + password + "'));";
    }

    let dataString = false;
    if (request.data) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString();
        }
        const parsedQueryString = querystring.parse(request.data, { sort: false });
        dataString = '$data = array(\n';
        const dataCount = Object.keys(parsedQueryString).length;
        if (dataCount === 1 && !parsedQueryString[Object.keys(parsedQueryString)[0]]) {
            dataString = "$data = '" + quote(request.data) + "';";
        } else {
            let dataIndex = 0;
            for (const key in parsedQueryString) {
                const value = parsedQueryString[key];
                dataString += "    '" + key + "' => '" + quote(value) + "'";
                if (dataIndex < dataCount - 1) {
                    dataString += ',\n';
                }
                dataIndex++;
            }
            dataString += '\n);';
        }
    }
    let requestLine = '$response = Requests::' + request.method + '(\'' + request.url + '\'';
    requestLine += ', $headers';
    if (dataString) {
        requestLine += ', $data';
    }
    if (optionsString) {
        requestLine += ', $options';
    }
    requestLine += ');';

    let phpCode = '<?php\n';
    phpCode += 'include(\'vendor/rmccue/requests/library/Requests.php\');\n';
    phpCode += 'Requests::register_autoloader();\n';
    phpCode += headerString + '\n';
    if (dataString) {
        phpCode += dataString + '\n';
    }
    if (optionsString) {
        phpCode += optionsString + '\n';
    }

    phpCode += requestLine;

    return phpCode + '\n';
}

const toCpr = (curlCommand) => {
    const request = parseCurlCommand(curlCommand);

    let code = '#include <iostream>\n';
    code += '#include <string>\n';
    code += '#include "cpr/cpr.h"\n';
    code += '#include "utils/helper.h"\n';
    code += '#include "nlohmann/json.hpp"\n\n';
    code += 'using namespace cpr;\n\n';
    code += 'int main(int argc, char **argv) {\n';
    code += '\tSession session;\n';

    if (request.cookieString) {
        code += `\tsession.SetCookies(common::Url::ParseCookies("${request.cookieString}"));\n`;
    }
    if (request.insecure) {
        code += '\tsession.SetVerifySsl(false)\n';
    }
    if (request.proxy) {
        request.proxy.indexOf('socks5') > -1 ? code += `\tsession.SetProxies({{"socks5", "${request.proxy}"}})` : code += `\tsession.SetProxies({{"http", "${request.proxy}"}, {"https", "${request.proxy}"}});\n`;
    }

    code += `\tstd::string url = "${request.url}";\n`;
    code += '\tsession.SetUrl(Url{url});\n';
    if (request.data) {
        if (typeof request.data === 'number') {
            request.data = request.data.toString();
        }
        code += `\tBody data = R"(${request.data})";\n`;
        code += '\tsession.SetBody(data);\n';
    }
    if (request.headers) {
        code += `\tHeader headers = common::Url::ParseHeaders(R"(${jsonIndent(request.headers, '\t\t')})");\n`;
        code += '\tsession.SetHeader(headers);\n';
    }
    code += '\tResponse r = session.' + capitalizeUpper(request.method) + '();\n';
    code += '\tstd::cout << r.text << std::endl;\n';

    code += '\treturn 0;\n';
    code += '}';

    return code + '\n';
}

const curlTransfer = (curlCommand, language, blank) => {
    var result;
    switch (language) {
        case "golang":
            result = toGolang(curlCommand);
            break;
        case "json":
            result = toJSONString(curlCommand);
            break;
        case "python":
            result = toPython(curlCommand);
            break;
        case "aiohttp":
            result = toPython(curlCommand, true);
            break;
        case "node-request":
            result = toNodeRequest(curlCommand);
            break;
        case "node-fetch":
            result = toNodeFetch(curlCommand);
            break;
        case "php":
            result = toPhp(curlCommand);
            break;
        case "cpp":
            result = toCpr(curlCommand);
            break;
        case "requests":
            result = toGoRequests(curlCommand);
            break;
        default:
            result = toGoRequests(curlCommand, "");
            break;
    }
    return result;
}

var args = process.argv.splice(2);
const ret = curlTransfer(args[0], args[1]);
console.log(ret);

const copyCommand = process.platform === "darwin" ? "pbcopy" : "clip";
exec(copyCommand, function (err, stdout, stderr) {
    if (err || stderr) return console.log(err, stdout, stderr);
}).stdin.end(ret);
