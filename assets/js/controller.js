var mTestApp = angular.module("mTestApp", ["ngAnimate", "ngSanitize", "ngCookies", "ngRoute", "ngTagsInput",
    "dndLists", "ui.bootstrap", "ngAnimate", "sticky", "switcher"]);

//----------------------------------------------------------------------------------------------------------
//---------------------------------------------COTROLLER WITH LOCAL ----------------------------------------
//----------------------------------------------------------------------------------------------------------

mTestApp.controller("mTestController", function($scope, $sce, $http, $cookies, $window,
    $location, LS, SubmitData, MathData, ModalWin, toPdf, $interval ) {
    $http.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
    //----------------------------------------------------------------------------------------------------------
    //---------------------------------------------DATA---------------------------------------------------------
    //----------------------------------------------------------------------------------------------------------

    this.value = LS.getData();
    this.data = angular.fromJson(this.value);
    $scope.disableSticking = false;
    this.num = MathData.getNum(this.data);

    $scope.emails = [];
    $scope.allowed = { dropzone: ['container'], container: ['itemplus', 'item'], itemplus: ['item'] };
    $scope.models = {
        selected: null,
        templates: [{
            type: "container",
            id: this.num + 1,
            columns: [
                []
            ],
            name: "Додати інф. вимогу"
        }, {
            type: "itemplus",
            id: 3,
            columns: [
                []
            ],
            name: "Додати складову інф. вимоги"
        }, {
            type: "item",
            id: 3,
            name: "Додати дію"
        }],
        dropzones: this.data
    };
    setInterval(function() {
        var lsData = angular.toJson($scope.models.dropzones);
        LS.setData(lsData);
    }, 5000);


    this.emails_arr = function() {
        var output = [];
        for (var k =0; k < $scope.emails.length; k++){
            output.push($scope.emails[k]['text'])
        }
        return output
    };


    $http({
        method: 'GET',
        url: '/static/1.json'
    }).then(function successCallback(response) {
         $scope.questions = response.data;
    }).catch(function (reason) { console.log(reason)});


    $scope.isLoggedIn = 234;

    this.mdata = function() {
        return angular.toJson($scope.models.dropzones);
    };

    $scope.getUrlData = SubmitData.getDataUrl;

    this.saveForSend = function($event) {
        $event.preventDefault
        SubmitData.EmailDataToDB(this.mdata(), $cookies.csrftoken)
    };
    this.multiActsSave = function($event){
        $event.preventDefault
        SubmitData.multiActsSend($scope.eauthor, $scope.eemail,
            this.mdata(), $cookies.csrftoken, this.emails_arr())
    }
    $scope.sendEmail = function($event) {
        $event.preventDefault
        return SubmitData.sendEmail($scope.eauthor, $scope.eemail,
            $scope.emails_arr(), $cookies.csrftoken, $scope.getUrlData())
    };

    $scope.sendEmailCorrupt = function($event) {
        $event.preventDefault
        return SubmitData.sendEmailCorr('M-TEST Corruption', 'mtest@clc.com.ua', 'vk@clc.com.ua',
            $cookies.csrftoken, toPdf.corruptEmail($scope.formData))
    };
    $scope.effectPopover = $sce.trustAsHtml('Критично – реалізація корупційного ризику призведе до зупинки діяльності <br> ' +
        'Суттєво - реалізація корупційного ризику призведе до суттєвих витрат грошей та/або часу співставних з розміром щомісячного прибутку');

    var trusted = {};
    $scope.getPopoverContent = function(content) {
        return trusted[content] || (trusted[content] = $sce.trustAsHtml(content));
    };

 //--------------------------------------------- ADD-REMOVE ROW ---------------------------------------------
   $scope.addRow = function(item,text){
    var quest = {
        text: text,
        "yesno": [{"name": "Так", "value": "yes"}, {"name": "Ні", "value": "no"}],
        "remove": "block"
    };

    item.push(quest);
    };

    $scope.removeRow = function(item, index){
        item.splice(index, 1);
   };
    $scope.deleteTrash = function(where, stndrt, real) {

        return where.slice(stndrt, real)
    };
    //---------------------------------------------SUM----------------------------------------------------------

    this.Sum = function(id) {
        return MathData.mathSum($scope.models.dropzones['1'], id)
    };
    this.totalSum = function() {
        return MathData.totalSum($scope.models.dropzones['1'][0]['ki'], this.Sum)
    };

    $scope.corr_sum = function(data){
        var total = 0;
        for (var k in data){
            total += Number(data[k])
        }
        return total;
    };
    $scope.addToSlider = function(res, eft){
        if (res == 1 && eft == 1) {
            return 1
        } else if (res == 1 && eft > 2){
            return 100
        } else {
            return 0
        }

    };
    $scope.getLength = function(obj) {
        if (obj != null) {
            return obj.length
        } else {
            setTimeout($scope.getLength, 300);
        }

    };
    //---------------------------------------------MODAL--------------------------------------------------------

    $scope.openModal = function(size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };

    $interval(function() {
        var a = document.getElementById('report').innerHTML;
        var blob = new Blob(["\ufeff", a], {
            type: 'text/html'
        });
        $scope.saveToPdf = (window.URL || window.webkitURL).createObjectURL(blob);
    }, 1000000);




    $scope.hideQuestion = function(dep, val) {
        if (dep == 1 && val == 0) {
            return "none"
        } else if (dep == 1 && val == 1) {
            return "block"
        } else {
            return "block"
        }
    };
    $scope.payContent = $sce.trustAsHtml('Оплата – зазвичай середня  вартість людино-години роботи персоналу, який\
     виконує дії необхідні для  виконання ІВ.<br>\
  За бажанням можна використовувати  заробітну плату однієї години відповідного Спеціаліста або Керівника. <br>\
  Для ФОПів  це може бути  розрахунок вартості  людино-години на основі річного обсягу реалізації сектору (детальніше див.  Посібник з М-Тесту).<br>\
  В віконце  оплата вводиться  за місяць і перераховується автоматично в оплату за годину.');


    //---------------------------------------------CORR_BAR--------------------------------------------------------
     $scope.hideCritical = function(val, length) {
        if (val >= 100 || val > length/2) {
            return "block"
        } else {
            return "none"
        }
    };

    $scope.corr_bar = function(total, length) {
    var a = '';
    if ( total > length/2){
        a = 'danger';
        $scope.barText = "Критична корупційна складова!";
    } else {
        $scope.barText = "Кількість корупційних ризиків:" + total;
        a = 'success'
    }
    return a
  };

    $scope.valueToText = function (val, type) {
        if (type == "yn") {
            if (val == 1 || val == "Yes"){
                return "Так"
            } else if (val == 0 || val == "No"){
                return "Ні"
            } else if (val == "idk"){
                return "Незнаю"
            } else {
                return "Не заповнено!"
            }
        } else {
            if (val == 1){
                return "Суттєво"
            } else if (val == 0){
                return "Несуттєво"
            } else if (val == 100){
                return "Критично"
            } else {
                return "Не заповнено!"
            }
        }
    };

    //---------------------------------------------RESET--------------------------------------------------------
    $scope.reset = function() {
        $scope.models.dropzones = angular.copy({"1":[{"type":"container","id":3,"columns":[[{"type":"itemplus","id":3,
                    "columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},{"type":"item","id":6,"name":"Додати дію","subsum":0}]],
                    "name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0},
                {"type":"container","id":null,"columns":[[{"type":"itemplus","id":4,"columns":[[{"type":"item","id":3,"name":"Додати дію","subsum":0},
                            {"type":"item","id":4,"name":"Додати дію","subsum":0}]],"name":"Додати складову інф. вимоги"}]],"name":"Додати інф. вимогу","contsub":0}]});
    }
    
});
mTestApp.controller("authLoginController", function ($scope, $timeout, $location, authService) {
    $scope.user = {};
    $scope.onLogin = function() {
        authService.login($scope.user)
            .then(function(user) {
                localStorage.setItem('token', user.data.token);
                $location.path('/u/cabinet');
            })
            .catch(function(err) {
                console.log(err);
                $scope.message = "Невірний логін або пароль, спробуйте ще раз";
                $timeout(function() {
                    $scope.message = ""
                }, 2000);
            });
    };
});

mTestApp.controller("userCabinetController", function ($scope, $http, $location, $rootScope,
                                                       $timeout, authService, mtCrud, ModalWin) {
    $scope.changepass = false
    const token = localStorage.getItem('token');
    if (token) {
        authService.ensureAuthenticated(token)
            .then(function(user) {
                if (user.status === 200) {
                    $scope.userdata = user.data.data;
                    $scope.records = user.data.data.records;
                    $rootScope.isLoggedIn = true;
                }
            })
            .catch(function(err) {
                console.log(err);
                $location.path('/u/login');
            });
    }
    $scope.changeUserField = function (field, id, value) {
        console.log(field, id, value);
        $http({
            method: 'POST',
            url:"/api/v.1/u/edituser",
            data: {field: field, data: value, id: parseInt(id)},
            headers: {'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function(err){
            console.log(err)
        });
    };
    $scope.openModal = function(size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };
    //--------------CRUD mtest------------------------
    $scope.addMtest = function (newmtest) {
        mtCrud.addMtest(newmtest, token).
        then(function (response) {
            console.log(response.data);
            $scope.records = response.data.records;
            $scope.len_records = angular.toJson($scope.records).length
        }).catch(function(err){
            console.log(err)
        });
    };
    $scope.removeMtestItem = function (id) {
        console.log(id);
        mtCrud.removeMtestItem(id, token)
            .then(function (response) {
            delete $scope.records[id];
            $scope.len_records = angular.toJson($scope.records).length
        }).catch(function(err){
            console.log(err)
        });
    };

    $scope.updateMtestItem = function (item) {
        console.log(item);
        mtCrud.updateMtestItem(item, token)
            .then(function (response) {
        }).catch(function(err){
            console.log(err)
        });
    };

    $scope.dynamicPopover = {
        isOpen: {},
        title: "heh",

        templateUrl: {},

        open: function open(index, template, tp) {
            $scope.dynamicPopover.templateUrl[index+tp] = template;
            $scope.dynamicPopover.isOpen[index + tp] = true;
            $scope.dynamicPopover.data = 'Hello!';
        },

        close: function close(index, tp) {
            $scope.dynamicPopover.isOpen[index + tp] = false;
        }
    };

    var countUp = function() {
        $scope.len_records = angular.toJson($scope.records).length
    };
    $timeout(countUp, 100);


});

mTestApp.controller("authRegisterController", function ($scope,authService) {

    $scope.user = {
        password: "",
        confirmPassword: ""
    };
    $scope.onRegister = function() {
        authService.register(vm.user)
            .then(function(response) {
                $location.path('/status');
            })
            .catch(function(err) {
                console.log(err);
            });
    };
});

mTestApp.controller("menuController", function ($scope, $rootScope, authService, ModalWin) {
    $rootScope.isLoggedIn = false;
    const token = localStorage.getItem('token');
    if (token) {
        authService.ensureAuthenticated(token)
            .then(function(user) {
                if (user.status === 200) {
                    $rootScope.isLoggedIn = true;
                }
            })
            .catch(function(err) {

                console.log(err)
            });

    }
    $scope.onLogout = function() {
        localStorage.removeItem('token');
        $rootScope.isLoggedIn = false;
        $location.path('/u/login');
    };
    $scope.openModal = function(size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };
});

mTestApp.controller("mTestDBController",
    function($scope, $timeout, $sce, $http, $cookies, $location, $routeParams, DB_data,
                                                   $window, LS, SubmitData, MathData, ModalWin, toPdf) {
    $scope.allowed = { dropzone: ['container'], container: ['itemplus', 'item'], itemplus: ['item'] };

    $scope.params = $routeParams
    //----------------------------------------------------------------------------------------------------------
    //---------------------------------------------DATA---------------------------------------------------------
    //----------------------------------------------------------------------------------------------------------
        SubmitData.getModelsData($scope.params.mtest_id)
            .then(function(response) {
                $scope.mtestData = response.data.mtest
                this.val = angular.fromJson(response.data.mtest.calculations)
                $scope.modelsdb = {
                    selected: null,
                    templates: [{
                        type: "container",
                        id: this.num + 1,
                        columns: [
                            [],
                        ],
                        name: "Додати інф. вимогу"
                    }, {
                        type: "itemplus",
                        id: 3,
                        columns: [
                            [],
                        ],
                        name: "Додати складову інф. вимоги"
                    }, {
                        type: "item",
                        id: 3,
                        name: "Додати дію"
                    }],
                    dropzones: this.val
                };
            });

    //-------------------------------------------- FOR TYPEHEAD ------------------------------------------------

    $http({
        method: 'GET',
        url: '/json/dii/dia_name'
    }).then(function(response) {
        $scope.dii = response.data;
    });



    //-------------------------------------------- GROUP MATH ------------------------------------------------

    setTimeout(function() {
        $scope.devs_calc = function () {

        }
        var devs = $scope.modelsdb.dropzones['1'][0].executors
        var infs = $scope.modelsdb.dropzones['1']

        for(var devs_get = [];devs_get.length < infs.length; devs_get.push([]));

        for(var i = 0; i < devs.length; i++){
            let dev = devs[i]
            let dev_data;

            if(dev[2] === true) {
                $http({method: 'GET', url: dev[1]+'/get/'}).then(function successCallback(response) {
                    dev_data = JSON.parse(response.data)
                })
                setTimeout(function() {
                    for (var j = 0; j < dev_data['1'].length; j++) {
                        if (dev_data['1'][j]['contsub'] != 0){
                            devs_get[j].push(dev_data['1'][j]['contsub'])
                        }
                    }
                }, 50)
            }

        }
        $scope.devs_get = devs_get

    }, 150);

    setTimeout(function() {
        $scope.sumDevsInfs = function(x) {
            var tmp_sum = 0;
            for (var i = 0; i < x.length; i++) {
                tmp_sum += x[i]
            }
            return tmp_sum/x.length
        }
    }, 200);

    $scope.reloadRoute = function() {
        $window.location.reload();
    }
    //-------------------------------------------- END GROUP MATH ------------------------------------------------


    //----------------------------------------------------------------------------------------------------------
    //---------------------------------------------SAVE TO DB --------------------------------------------------
    //----------------------------------------------------------------------------------------------------------

    this.mdata = function() {
        return angular.toJson($scope.modelsdb.dropzones);
    };
    this.submitData = function($event) {
        $event.preventDefault
        SubmitData.saveDataToDB(this.mdata(), $cookies.csrftoken)
    };

    //---------------------------------------------SUM----------------------------------------------------------

    $scope.SumDb = function(id) {
        return MathData.mathSum($scope.modelsdb.dropzones['1'], id)
    };
    $scope.totalSumDb = function() {
        return MathData.totalSum($scope.modelsdb.dropzones['1'][0]['ki'], $scope.SumDb)
    };

    $scope.awgSumDb = function() {
        return MathData.awgSum($scope.modelsdb.dropzones['1'], $scope.modelsdb.dropzones['1'][0]['ki'])
    };

    //---------------------------------------------MODAL--------------------------------------------------------
    this.openModal = function(size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };



    //---------------------------------------------alert------------------------------------------------------


    $scope.getClass = 'not_disabled'
    $scope.getClassTwo = 'disabled'

    this.setClasses = function() {
        $scope.getClass = 'disabled'
        $scope.getClassTwo = 'not_disabled'
        $timeout(function() {
            $scope.getClass = 'not_disabled', $scope.getClassTwo = 'disabled'
        }, 1000)
    }


    //---------------------------------------------PDF--------------------------------------------------------
    //$scope.$watch('modelsdb.dropzones', function(dropzones) {
    //    if (dropzones) {
    //        var data = toPdf.saveAsHtml($scope.modelsdb.dropzones['1'], $scope.SumDb(), $scope.totalSumDb())
    //        var blob = new Blob(["\ufeff", data], {
    //            type: 'text/html'
    //        });
    //        $scope.saveToPdf = (window.URL || window.webkitURL).createObjectURL(blob);
    //    }
    //});
//---------------------------------------------reset------------------------------------------------------
    $scope.resetdb = function() {
        SubmitData.getModelsData().then(function(response) {
            val = angular.fromJson(response.data)
            $scope.modelsdb.dropzones = angular.copy(
                val
            )
        })
    };


});