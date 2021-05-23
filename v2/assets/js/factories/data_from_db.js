//----------------------------------------------------------------------------------------------------------
//--------------------------------------------- DB FACTORY -------------------------------------------
//----------------------------------------------------------------------------------------------------------
var mTestApp = angular.module('mTestApp');
mTestApp.factory('DB_data', function ($http, $location) {
    var data;
    return {
        getData: function (callback) {
            if(data) {
                callback(data);
            } else {
                $http.get($location.absUrl() + '/get/').success(function(d) {
                    callback(data = d);
                });
            }
        }
    };
});