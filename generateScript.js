const cookie = require('cookie');
const yargs = require('yargs');
const URL = require('url');
const querystring = require('query-string');

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

    let cookieString;
    let cookies;
    let url = parsedArguments._[1];

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

const curl2GoRequests = curlCommand => {
    const request = parseCurlCommand(curlCommand);
    let code = 'package main\n\n';
    code += 'import (\n\t"github.com/Esbiya/requests"\n\t"log"\n)\n\n';
    code += 'func main() {\n';
    code += `\turl := ${request.url.replaceAll("'", '"')}\n`;
    let args = "requests.RequestArgs{\n";
    if (request.headers) {
        code += `\theaders := requests.DataMap${jsonIndent(request.headers, '\t\t')}\n`;
        args += "\t\tHeaders: headers,\n";
    }
    if (request.cookies) {
        code += `\tcookies := requests.DataMap${jsonIndent(request.cookies, '\t\t')}\n`;
        args += "\t\tCookies: cookies,\n";
    }
    if (request.auth) {
        const splitAuth = request.auth.split(':');
        const user = splitAuth[0] || '';
        const password = splitAuth[1] || '';
        code += `\tauth := requests.DataMap${jsonIndent({username: user, password: password}, '\t\t')}\n`;
        args += "\t\tAuth:    auth,\n";
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
            args += "\t\tJSON:    data,\n";
            code += `\tdata := requests.DataMap${jsonIndent(data, '\t\t')}\n`;
        } catch (e) {
            if (request.data.indexOf("&") > -1) {
                data = JSON.parse(str2json(request.data));
                args += "\t\tData:    data,\n";
                code += `\tdata := requests.DataMap${jsonIndent(data, '\t\t')}\n`;
            } else {
                args += "\t\tBinary:  data,\n";
                code += `\tdata := []byte("${request.data}")\n`
            }
        }
    }
    if (!request.compressed) {
        args += "\t\tDisableCompression: true,\n";
        args = padSpace(args, "Headers", 11);
        args = padSpace(args, "Cookies", 11);
        args = padSpace(args, "Auth", 11);
        args = padSpace(args, "Data", 11);
        args = padSpace(args, "JSON", 11);
    }
    code += '\tresp := requests.' + capitalizeUpper(request.method) + `(url, ${args}\t})\n`;
    code += '\tif resp.Err != nil {\n\t\tlog.Fatal(resp.Err)\n\t}\n';
    code += '\tlog.Println(resp.Text)\n';
    code += '}';

    return code + '\n';
};

var args = process.argv.splice(2);
console.log(curl2GoRequests(args[0]));
