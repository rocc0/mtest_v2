mTestApp.service('authService', function ($http) {
    /*jshint validthis: true */
    const baseURL = 'http://localhost:8889/';

    this.login = function(user) {
        return $http({
            method: 'POST',
            url: baseURL + 'u/login',
            data: user,
            headers: {'Content-Type': 'application/json'}
        });
    };
    this.logout = function(user) {
        return $http({
            method: 'POST',
            url: baseURL + 'u/logout',
            data: user,
            headers: {'Content-Type': 'application/json'}
        });
    };
    this.register =  function(user) {
        return $http({
            method: 'POST',
            url: baseURL + 'u/register',
            data: user,
            headers: {'Content-Type': 'application/json'}
        });
    };
    this.resetpass =  function(email) {
        return $http({
            method: 'POST',
            url: baseURL + 'api/v.1/u/reset/',
            data: {email: email},
            headers: {'Content-Type': 'application/json'}
        });
    };
    this.setnewpass = function (pass, hash) {
        return $http({
            method: 'POST',
            url: baseURL + 'api/v.1/u/reset/' + hash,
            data: {password: pass}
        })
    }
    this.checkhash = function (hash) {
        return $http({
            method: 'GET',
            url: baseURL + 'api/v.1/u/reset/' + hash
        })
    }
    this.ensureAuthenticated = function(token) {
        return $http({
            method: 'GET',
            url: baseURL + 'api/v.1/u/cabinet',
            headers: {
                'Content-Type': 'application/json',
                Authorization: 'Bearer ' + token
            }
        });
    };
});

mTestApp.service('mtCrud', function ($http) {
    const baseURL = 'http://localhost:8889';

    this.readMtestFromDB = function (id) {
        return $http({
            method: 'GET',
            url: '/api/v.1/m/get/'+id
        });
    };

    this.addMtest = function(newmtest, token) {
        var data = {
            name: newmtest.name,
            region: parseInt(newmtest.region.id),
            government: parseInt(newmtest.government.id),
            calc_type: parseInt(newmtest.calc_type)
        };
        if (newmtest.calc_type === 1) {
            data.executors = {}
        }
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/create",
            data: data,
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
    this.removeMtestItem = function (id, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/delete",
            data: {id:id},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };

    this.updateMtestItem = function (item, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/update",
            data: {mid:item.id, name:item.name,
                region: parseInt(item.region),govern: parseInt(item.govern),
                calc_type: parseInt(item.calc_type), executors: item.executors },
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
    this.updateMtestCalculations= function (id, item, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/update",
            data: {id:id,calculations: item},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
    this.updateMtestExecutors= function (id, item, token) {
        console.log(item);
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/update",
            data: {id:id, executors: item},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
    this.addMtestExecutor = function (title, email, region, gov, dev_mid, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/excreate",
            data: {title: title, email: email, region: region, government: gov, dev_mid: dev_mid, },
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    };
    this.removeMtestExecutor = function (email, exIndex, devIndex, token) {
        return $http({
            method: 'POST',
            url:baseURL + "/api/v.1/m/exdelete",
            data: {ex_email: email, ex_mtest_id: exIndex, dev_mtest_id: devIndex},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        });
    }
});

mTestApp.filter('objLength', function() {
    return function(object) {
        var count = 0;

        for(var i in object){
            count++;
        }
        return count;
    }
});