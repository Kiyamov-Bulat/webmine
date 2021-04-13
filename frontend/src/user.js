const ANONYMOUS = 0, LOGGED = 1

// User's group
const USER      = 0, 
      ADMIN     = 1,
	    MODERATOR = 2,
	    EDITOR    = 3

export function User(userObj) {
    userObj.isLogged = function() {
        return typeof userObj.id === 'number' && userObj.id != 0
    }
    return userObj
}

export function getUserFromCookie() {
    try {
        const json = getCookie("user")
        return new User(JSON.parse(json))
    } catch(e) {
        if (e.name == "SyntaxError") {
            return new User({})
        } else {
            throw e
        }
    }
}

export function saveUserInCookie(user) {
    let date = new Date(Date.now() + 86400e3);
    setCookie("user", JSON.stringify(user), {expires: date})
}

function getCookie(name) {
    let matches = document.cookie.match(new RegExp(
      "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
    ));
    return matches ? decodeURIComponent(matches[1]) : undefined;
  }

function setCookie(name, value, options = {}) {

    options = {
      path: '/',
      ...options // при необходимости добавьте другие значения по умолчанию
    };
  
    if (options.expires instanceof Date) {
      options.expires = options.expires.toUTCString();
    }
  
    let updatedCookie = encodeURIComponent(name) + "=" + encodeURIComponent(value);
  
    for (let optionKey in options) {
      updatedCookie += "; " + optionKey;
      let optionValue = options[optionKey];
      if (optionValue !== true) {
        updatedCookie += "=" + optionValue;
      }
    }
  
    document.cookie = updatedCookie;
  }

  function deleteCookie(name) {
    setCookie(name, "", {
      'max-age': -1
    })
  }