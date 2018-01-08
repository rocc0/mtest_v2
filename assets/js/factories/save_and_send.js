var mTestApp = angular.module('mTestApp');
mTestApp.factory('SubmitData', function ($window, $http, $location,$cookies) {
    var urlData = '';
    return {
        saveDataToDB: function (mdata, token) {
            var url = $location.absUrl();
            var data = $.param({ author: 'Admin', content: mdata, csrfmiddlewaretoken: token });
            $http.post(url, data).then(function (response) {
                console.log(response.headers('Redirect-URL'))
                })
        },
        EmailDataToDB: function (mdata, token) {
            var url = $location.absUrl();
            var data = $.param({ author: 'Admin', content: mdata, csrfmiddlewaretoken: token });
            $http.post(url, data).then(function (response) {
                urlData = response.headers('Redirect-URL')
                return urlData
            })
        },
        sendEmail: function (author, email, emailfor, token, message, redirect) {
            var data = $.param({
                csrfmiddlewaretoken: token,
                author: author,
                emailfrom: email,
                emailfor: emailfor,
                message: message});
            $http.post("/email/send/", data).then(function (response) {
                console.log(response.headers)
                $window.location.href = message[0]
            })
        },
        multiActsSend: function (author, email, mdata, token, addresses) {
            var url = $location.absUrl();
            var acts = [];
            if(addresses.indexOf(email) === -1){
                addresses.push(email)
            }
            for (i = 0; i < addresses.length; i++) {
                let executor = addresses[i];
                let tmp_data = angular.fromJson(mdata)
                tmp_data['1'][0].rigths = 'no'
                mdata = angular.toJson(tmp_data)
                var data = $.param({ author: executor, content: mdata, csrfmiddlewaretoken: token });
                $http.post(url, data).then(function (response) {
                    let redir = '/e/' + response.headers('Redirect-URL')
                    acts.push([executor, redir])
                    if (email != executor){
                        var emaildata = $.param({csrfmiddlewaretoken: token, author: author, emailfrom: email,
                        emailfor: executor, message: 'http://192.168.99.1:8000/e/' + response.headers('Redirect-URL')});
                    $http.post("/email/send/", emaildata).then(function (response) {
                    })
                    }

                    return acts
                })
            }
            setTimeout(function() {
                var tmp_data = angular.fromJson(mdata)
                tmp_data['1'][0].rigths = 'yes'
                tmp_data['1'][0].executors = acts
                mdata = angular.toJson(tmp_data)
                var devdata = $.param({ author: author, content: mdata, csrfmiddlewaretoken: token });
                var devurl = 'http://192.168.99.1:8000/d/' + acts[acts.length - 1][1].slice(3,acts[acts.length - 1][1].length)
                $http.post(devurl, devdata).then(function (response){
                let devurl = 'http://192.168.99.1:8000/d/' + response.headers('Redirect-URL')
                var emaildata = $.param({csrfmiddlewaretoken: token, author: 'Mtest',
                    emailfrom: email, emailfor: email, message: 'dev'+'|'+acts.slice(0,acts.length -1).join(';')+'|'+devurl});
                $http.post("/email/send/", emaildata).then(function (response) {
                    $window.location.href = devurl
                })
            })
            }, 1000)


            return acts
        },
        getDataUrl: function () {
            var currentUrl = $location.absUrl()
            if (currentUrl == 'http://192.168.99.1:8000/') {
                return 'http://192.168.99.1:8000'+urlData
            } else {
                return currentUrl
            }
        },
        getModelsData: function (id) {
            console.log(id);
            return $http({
                method: 'GET',
                url: '/api/v.1/m/get/'+id
            });
        },
        sendEmailCorr: function (author, email, emailfor, token, message) {
            var data = $.param({
            csrfmiddlewaretoken: token,
            author: author,
            emailfrom: email,
            emailfor: emailfor,
            message: message
            });
        $http.post("/email/send/", data).then(function (response) {
            console.log(response.headers)
            });
        }
    };
});