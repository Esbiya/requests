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
            if (typeof parsedArguments[argName] === 'string') {
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
            const value = splitArguments[1];
            multipartUploads[key] = value;
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
        args += "Headers: headers,\n";
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
        if (request.data.indexOf("&") > -1) {
            request.data = str2json(request.data);
            args += "\t\tData:    data,\n";
        } else {
            args += "\t\tJSON:    data,\n";
        }
        code += `\tdata := requests.DataMap${jsonIndent(JSON.parse(request.data), '\t\t')}\n`;
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

console.log(curl2GoRequests(`curl 'https://login.taobao.com/newlogin/login.do?appName=taobao&fromSite=0&_bx-v=1.1.20' \\
  -H 'authority: login.taobao.com' \\
  -H 'pragma: no-cache' \\
  -H 'cache-control: no-cache' \\
  -H 'eagleeye-sessionid: t6kgpjLFz2m58n10zk21xtIrgabF' \\
  -H 'accept: application/json, text/plain, */*' \\
  -H 'eagleeye-pappname: gf3el0xc6g@256d85bbd150cf1' \\
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36' \\
  -H 'eagleeye-traceid: b216fbdc1610766945306100250cf1' \\
  -H 'content-type: application/x-www-form-urlencoded' \\
  -H 'origin: https://login.taobao.com' \\
  -H 'sec-fetch-site: same-origin' \\
  -H 'sec-fetch-mode: cors' \\
  -H 'sec-fetch-dest: empty' \\
  -H 'referer: https://login.taobao.com/member/login.jhtml?style=miniall&newMini2=true&full_redirect=true&&redirectURL=http%3A%2F%2Fworld.taobao.com%2F&from=worldlogin&minipara=1,1,1&from=worldlogin' \\
  -H 'accept-language: zh-CN,zh;q=0.9' \\
  -H 'cookie: tracknick=; miid=1928852900179986723; sgcookie=E1UIaaaplkSrC4cTJf4kP; _samesite_flag_=true; xlly_s=1; _m_h5_tk=7c054d74754411dc668bce986dd75216_1610775569864; _m_h5_tk_enc=55a0fdf4a40433f9ae9d0594c3c4f42e; cookie2=1c386caaf67de4d7690e183d5f775c4d; t=4e7cb12c7ff6bcb67d8b026c2a526d12; _tb_token_=e1d14f57e6ee3; cna=xC68F5h8QW8CATsq7e47o7pt; _bl_uid=88kaOj6wzI55Us131m251mUu1tCC; mt=ci=-1_0; _fbp=fb.1.1610766936473.1947231865; uc1=cookie14=Uoe1gq9r2Iujww%3D%3D; hng=CN%7Czh-CN%7CCNY%7C156; thw=cn; XSRF-TOKEN=a317e807-7ec1-4262-951b-5a24d963b790; tfstk=c7QdBombV5hKHLOOgMEiFyRrF-NcZdFpQXOKySbXJXHbW6uRiR7ckkg0ApkpW5C..; l=eBOYAwmHQG1jr4gvBOfZnurza779hIRcguPzaNbMiOCPO7fH5xMAW6G4BHLMCnGVnsEBr3zoKBo0B-LKZyUgl6Yl3ZQ7XPQoPdTh.; isg=BGpqw48L_idKrkx8cro3jJM0u9YM2-416jfPyfQjdL1IJwnh1muRRVGRt1M712bN' \\
  --data-raw 'loginId=18829040039&password2=84297912adeed3f348b3da7c2df2cec540f37423d748f3323f8e1e183848e41d63fb5f5bb222a3e333370c0c454fc364b06ab588c9d9a4ece3879d95f30d62ea9f87384edf16242502e457cbe52b12e7d7afce0b45033ff2cc1864ed40afccdcfd3eff46efb3d933fdc5f1d00872458ee9586c2dfb4fbfd7ec36f2cceab3fb6b&keepLogin=false&ua=140%23PrXDuQbczzFxwQo22Z3TCtSdvWMkKWzgkCqlMM6vlUV6B1j5G5TbbT8jtnpqJRQFCOSTi2n4HZSaCE3n5gycrvWqlbzx6IUpctgqzzr2hX0iU6OzzPzbVXlqlbrrNwQ1txyWEAziL2ILllfzzPzibvjj0TTd2eDKQt8NzeOVu2Ygl%2BFozD3RrfV9ONdOHaU3%2F0ttkrcT75HE391Y8Wc6THI4ygA%2FGlJe6QjYf72oLd%2FADbHH3r68kHTQcX946%2BVsEDjAV0MyefSu6ysv0r8l4ix%2FqAQj6dBt7oPHb%2FGFNL2xp1M1EzChs6OUQ47qVc2k79rvGNvkQs8aC5%2FvzajwMx5lOuxMn9oJmRPG37RE3uI2POJJf03NiElIZbXBL363wnO8pkz%2BaPrJoGvJvhI2acEFW5t8y5hs146Cz1D0ZBp1Rm%2BoNy%2FHZfP9ivfUi1vBV3nAUNsS78ex03If27zIQxJBcRlbq7ak%2Bxgk%2BsESuyx9VDJo8MMs7zywECZlaKr6E02gHJPbOHlFwHzPoHtYuGxYxOFpw9eNQ6g000GGMMQqchBa23%2BWLbrOSb9hArGNLfnxkq0n%2FRVD8QHnZ8xL2KwqOzjYQUmgdUkgJY40V7jkbU7mqFmIjqJWW6psGyi7Q1X8Hh8oubL2DHvxFR292VTG2x6CNB1mFMJYF8MdWIHMZVVZ8g%2FerRXG%2FXqMz5ftDkYSevb3t%2BBcgoKlRTx6nSc6vzYRu4yJW%2BjRMN0T8EO0GtrhGFgFkrb6Ms9D2HA2%2F23YVnCGaGyLag3T5bmzLJzx82r2eijrvP0R3def8mGVMQuWO17pdwpeepGUFu6VThegt6A8KvPqAWwRMSl9kqEaThZwawiXHHRkSVejwIgRiRrEml9zdsBvWZ%2FVPNAzsSrnUBMCxbjZSI4sX6nX0424o9%2BSlUB%2BoCPJcud%2Be6u4okn2WQvvn47BGAtKac1Hzv4LZE4ZwZa%2FunB61%2B7NWAqZ43ZEmVV%2FUWsH40K2jOgEyeHBoQ9PlecJB46sEPMlDL4nRgefk1QWK7DptICl6dxHRk1hsGM%2FOjbnIW4v9q3ILOuWbZI6YRKV93idOPGffaEEBmDWnyj2psMWvrasTKikmUQz1G2J0JJY3z6OWvL2kGV%2FnljF0TTdH7bofSNxLk0FbWAi6ZWQUaD%2FAP%2B%3D&umidGetStatusVal=255&screenPixel=2560x1440&navlanguage=zh-CN&navUserAgent=Mozilla%2F5.0%20%28Macintosh%3B%20Intel%20Mac%20OS%20X%2010_15_7%29%20AppleWebKit%2F537.36%20%28KHTML%2C%20like%20Gecko%29%20Chrome%2F87.0.4280.141%20Safari%2F537.36&navPlatform=MacIntel&appName=taobao&minipara=1%2C1%2C1&appEntrance=taobao_pc&_csrf_token=1WhXEGVzpZkfPhONEGo3K6&umidToken=71fbc9edf087812e8959453cc22590967f38bdeb&hsiz=1c386caaf67de4d7690e183d5f775c4d&newMini2=true&bizParams=&full_redirect=true&style=miniall&appkey=00000000&from=worldlogin&isMobile=false&lang=zh_CN&returnUrl=http%3A%2F%2Fworld.taobao.com%2F&fromSite=0&bx-ua=140#paur3bbqzzF0rzo22Z3TCtSdvWMkKWzgkCqlMM6vlUV6B1j5G5TbbT8jtnpqJRQFCOSTi2n4HZSaCE3n5gy3d8gqlbzx6IUpctgqzzr141HaU6OzzPzbVXlqlbrdNwQ1txyWEAD2Q282UpszaIziVXEFLrfxh8wKJp8Wzl2V8OYNlKMocD+bVJMoLhGP7SvZrI7ZbixvdsbhJ3+tNWB0D7r5Q/6lVY/dLqurkascPcmmze9R9ZssgbxjzFchXZyGe2gFBu4J0WLjXpfr8CbdUdQJXfWMLKY55RVXl/UeJq6p2DA22WdGYDyGKSQJjsGZuCB+Uuyap8ZHqroF9tHFdvtcKqkHBdtasn/5bxx1UryAJRcF2soPguernSVljOA2sEMjNYLgunkj0PoTcy4It0B7NmQsj/4eudPbhIcur0WDqxLv20MtMo59iZexGU3CcXa0m5/+AG5KrjDNUdCoOcgpxKUWuXJfAh6csfSQ3LJn3EADNn9U/VUm56YibY29LAQgKzF9eCOYp39ViX5OThNZiK8prYfoE7nf5aGz3tdS3Bi3mi3S+RZBP6XBZuAxMTBTaQWgZ6uUQ47qVc2k79rvGNvkQs8aC5/vzajwMx5lOuxMnjqnzZPDilBOYyL/NVMGxFLKAj1cuYv7mW8on6gr+Hv+Yj+VFI2juIx1uv4DUX5S2rqnmkPTiD8361zJZw13J55vMEzOxenMbbtkyCVWz4yXsPT7e+Q2B7xdIjGh9RUFYpSd0Ljj/YWVZi87NShkabt2Rf+T2jErB+EpFpSuiQuBnUXQiVTjXY58E0iGaNIm6KbMHX3jfdg2Ia9VYEbEaosrwoCjX/jYgESIH8+jxnppD24cpCgxESTJ7zUupRe1VfuyEwVebVQf4g6Ke82qqCIz+KZrtbE2rnSXURWXgHuwt21S9JAv7wJlqYGfrCbClp46R84caETrBWwqM/YWNcjyeKDsbbEX0lSvZZ1nDAFMdvA5EINAZi/Bw9BrRpil1IfXpTIqXX/4gSCCGtq7rRpUAd0+K40SmcDRz1wOzBf+XH3eHlTDy07a9vrxCzfYS6e7wPD+g9bcryhZaRduvph4SIqSDaieQf+ipV6UaVtmnDPO77X1fgs0yQOVHaupfn7rq8K1DIYm0cXsoNqSshFLbouN0dXDNggYU4b6QYhXbO5lf0qCPyxl7qo0BGev7psZEIvHmTzqVc6YDAfJlIeLqu1Uu9Z6NRw5h+So6dqUCtG4jB8NlCH2IX2RJosoIJLn1TyLYKDpLJwTwcyn7hyP4ilvkhST18Ht4scPNtzR9GhKfSkESqYL64oGtmkXO74Lrr//LakMfY+IiTeOSAVBGgT39vmMzeGrlmVhTadZ3jd0wj212asJ4cR3oR8+bNwntb34qlrTd7Zp7O6hGFZeTwc8jacRBJ6/aznqcF==&bx-umidtoken=T2gAebgRYOzMeGlZ1saou1UjOJaT5_HsJJmNn9LqnShe53XQwi_6wCcC8EiqWka7Cwg=' \\
  --compressed`));
