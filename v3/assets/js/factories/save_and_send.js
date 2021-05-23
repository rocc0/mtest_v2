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
        getModelsData: function (id) {
            console.log(id);
            return $http({
                method: 'GET',
                url: '/api/v.1/m/get/'+id
            });
        }
    };
});