var mTestApp = angular.module("mTestApp", ["ngAnimate", "ngSanitize", "ngCookies", "ngRoute", "ngTagsInput",
    "dndLists", "ui.bootstrap", "ngAnimate", "sticky", "switcher"]);

//----------------------------------------------------------------------------------------------------------
//---------------------------------------------CONTROLLER WITH LOCAL ----------------------------------------
//----------------------------------------------------------------------------------------------------------

mTestApp.controller("mTestController", function ($scope, $rootScope, $sce, $http, $cookies, $window,
                                                 $location, LS, SubmitData, MathData, ModalWin, $interval) {
    $http.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
    //----------------------------------------------------------------------------------------------------------
    //---------------------------------------------DATA---------------------------------------------------------
    //----------------------------------------------------------------------------------------------------------

    this.value = LS.getData();
    this.data = angular.fromJson(this.value);
    $scope.disableSticking = false;
    this.num = MathData.getNum(this.data);

    $scope.allowed = {dropzone: ['container'], container: ['itemplus', 'item'], itemplus: ['item']};
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

    setInterval(function () {
        var lsData = angular.toJson($scope.models.dropzones);
        LS.setData(lsData);
    }, 5000);

    $http({
        method: 'GET',
        url: '/api/v.1/actions'
    }).then(function (response) {
        $scope.adm_actions = response.data.actions;
    });


    $http({
        method: 'GET',
        url: '/static/questions.json'
    }).then(function (response) {
        $scope.questions = response.data;
    }).catch(function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/static/help.json'
    }).then(function (response) {
        $rootScope.help_texts = response.data;
    }).catch(function (reason) {
        console.log(reason)
    });

    $scope.isLoggedIn = 234;

    $scope.mdata = function () {
        return angular.toJson($scope.models.dropzones);
    };

    $scope.effectPopover = $sce.trustAsHtml('Критично – реалізація корупційного ризику призведе до зупинки діяльності <br> ' +
        'Суттєво - реалізація корупційного ризику призведе до суттєвих витрат грошей та/або часу співставних з розміром щомісячного прибутку');

    var trusted = {};

    $scope.getPopoverContent = function (content) {
        return trusted[content] || (trusted[content] = $sce.trustAsHtml(content));
    };

    //--------------------------------------------- ADD-REMOVE ROW ---------------------------------------------
    $scope.addRow = function (item, text) {
        var quest = {
            text: text,
            "yesno": [{"name": "Так", "value": "yes"}, {"name": "Ні", "value": "no"}],
            "remove": "block"
        };
        item.push(quest);
    };

    $scope.removeRow = function (item, index) {
        item.splice(index, 1);
    };

    $scope.deleteTrash = function (where, stndrt, real) {
        return where.slice(stndrt, real)
    };
    //---------------------------------------------SUM----------------------------------------------------------

    $scope.Sum = function (id) {
        return MathData.mathSum($scope.models.dropzones['1'], id)
    };

    $scope.totalSum = function () {
        return MathData.totalSum($scope.models.dropzones['1'][0]['ki'], this.Sum)
    };

    $scope.corr_sum = function (data) {
        var total = 0;
        for (var k in data) {
            total += Number(data[k])
        }
        return total;
    };

    $scope.addToSlider = function (res, eft) {
        if (res == 1 && eft == 1) {
            return 1
        } else if (res == 1 && eft > 2) {
            return 100
        } else {
            return 0
        }

    };

    $scope.getLength = function (obj) {
        if (obj != null) {
            return obj.length
        } else {
            setTimeout($scope.getLength, 300);
        }

    };
    //---------------------------------------------MODAL--------------------------------------------------------

    $scope.openModal = function (size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };



    $scope.hideQuestion = function (dep, val) {
        if (dep === 1 && val === 0) {
            return "none"
        } else if (dep === 1 && val === 1) {
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
    $scope.hideCritical = function (val, length) {
        if (val >= 100 || val > length / 2) {
            return "block"
        } else {
            return "none"
        }
    };

    $scope.corr_bar = function (total, length) {
        if (total > length / 2) {
            $scope.barText = "Критична корупційна складова!";
            return 1;
        } else {
            $scope.barText = "Кількість корупційних ризиків:" + total;
            return 0;
        }
    };

    $scope.valueToText = function (val, type) {
        if (type == "yn") {
            if (val == 1 || val == "Yes") {
                return "Так"
            } else if (val == 0 || val == "No") {
                return "Ні"
            } else if (val == "idk") {
                return "Незнаю"
            } else {
                return "Не заповнено!"
            }
        } else {
            if (val == 1) {
                return "Суттєво"
            } else if (val == 0) {
                return "Несуттєво"
            } else if (val == 100) {
                return "Критично"
            } else {
                console.log(val)
                return "Не заповнено!"
            }
        }
    };

    $interval(function() {
        var a = document.getElementById('report').innerHTML;
        var blob = new Blob(["\ufeff", a], {
            type: 'text/html'
        });
        $scope.saveToPdf = (window.URL || window.webkitURL).createObjectURL(blob);
    }, 1000);

    //---------------------------------------------RESET--------------------------------------------------------
    $scope.reset = function () {
        $scope.models.dropzones = angular.copy({
            "1": [{
                "type": "container", "id": 3, "columns": [[{
                    "type": "itemplus", "id": 3,
                    "columns": [[{"type": "item", "id": 3, "name": "Додати дію", "subsum": 0}, {
                        "type": "item",
                        "id": 6,
                        "name": "Додати дію",
                        "subsum": 0
                    }]],
                    "name": "Додати складову інф. вимоги"
                }]], "name": "Додати інф. вимогу", "contsub": 0
            },
                {
                    "type": "container", "id": null, "columns": [[{
                        "type": "itemplus",
                        "id": 4,
                        "columns": [[{"type": "item", "id": 3, "name": "Додати дію", "subsum": 0},
                            {"type": "item", "id": 4, "name": "Додати дію", "subsum": 0}]],
                        "name": "Додати складову інф. вимоги"
                    }]], "name": "Додати інф. вимогу", "contsub": 0
                }]
        });
    }
    // popover for setting, mail, add mtest etc..
    $rootScope.seq = {
        "chat_bot": "registration",
        "registration": "regact",
        "regact": "item",
        "item": "dnd_button",
        "dnd_button": "inf_req",
        "inf_req":"inf_req_part","inf_req_part":"move",
        "move":"time","time":"payment", "payment":"usage_freq",
        "usage_freq":"additional_payment","additional_payment":"end"
    }


    $rootScope.calc_steps = [
        "inf_req","inf_req_part",
         "move",  "time",
        "payment","usage_freq",
        "additional_payment",
        "req_price", "end"
    ]

    $rootScope.contains = function (step) {
        for (let i = 0; i < $rootScope.calc_steps.length; i++) {
            if (step == $rootScope.calc_steps[i]) {
                return true
            }
        }
        return false
    }

    $rootScope.helpPopover = {
        showHelp: false,
        isOpen: {},
        title: "heh",
        content: "",
        lastStep: "",
        lastOpened: "",
        lastItem: 0,
        lastItemPlus: 0,
        lastMove: 0,
        open: function open(step) {
            $rootScope.helpPopover.templateUrl = "static/html/tmpl/modal_help_step.html";
            $rootScope.helpPopover.isOpen[step] = $rootScope.helpPopover.showHelp;
            $rootScope.helpPopover.content = $rootScope.help_texts[step];
            $rootScope.helpPopover.lastStep = step;
            $rootScope.helpPopover.lastOpened = step;
        },

        next: function next() {
            //cancel previous step
            $rootScope.helpPopover.isOpen[$rootScope.helpPopover.lastOpened] = false
            //get next index
            var ind = "";
            if ($rootScope.contains($rootScope.helpPopover.lastStep) || ($rootScope.helpPopover.lastStep == "dnd_button")) {
                if ($rootScope.seq[$rootScope.helpPopover.lastStep] == "end" || $rootScope.helpPopover.lastStep == "req_price") {
                    //find next valid item
                    if ($scope.models.dropzones[1][$rootScope.helpPopover.lastItem].columns[0][$rootScope.helpPopover.lastItemPlus].columns[0][$rootScope.helpPopover.lastMove+1] != undefined) {
                        $rootScope.helpPopover.lastMove++
                        $rootScope.helpPopover.lastStep = "inf_req_part"
                        ind = $rootScope.helpPopover.lastItem.toString()+"|"+
                            $rootScope.helpPopover.lastItemPlus.toString()+"|"+
                            $rootScope.helpPopover.lastMove.toString()+$rootScope.seq[$rootScope.helpPopover.lastStep];
                    } else if ($scope.models.dropzones[1][$rootScope.helpPopover.lastItem].columns[0][$rootScope.helpPopover.lastItemPlus+1] != undefined) {
                        $rootScope.helpPopover.lastItemPlus++;
                        $rootScope.helpPopover.lastStep = "inf_req"
                        ind = $rootScope.helpPopover.lastItem.toString()+"|"+
                            $rootScope.helpPopover.lastItemPlus.toString()+
                            $rootScope.seq[$rootScope.helpPopover.lastStep];
                        $rootScope.helpPopover.lastMove = 0;
                    } else if ($scope.models.dropzones[1][$rootScope.helpPopover.lastItem+1] != undefined) {
                        if ($rootScope.helpPopover.lastStep == "req_price") {
                            $rootScope.helpPopover.lastItem++;
                            $rootScope.helpPopover.lastStep = "dnd_button"
                            ind = $rootScope.helpPopover.lastItem.toString()+$rootScope.seq[$rootScope.helpPopover.lastStep];
                            $rootScope.helpPopover.lastItemPlus = 0;
                            $rootScope.helpPopover.lastMove = 0;
                        } else {
                            $rootScope.helpPopover.lastStep = "req_price"
                            ind = $rootScope.helpPopover.lastItem.toString()+"req_price"
                            $rootScope.helpPopover.isOpen[ind] = true
                            $rootScope.helpPopover.content = $rootScope.help_texts[$rootScope.helpPopover.lastStep];
                            $rootScope.helpPopover.lastOpened = ind
                            return;
                        }
                    } else if ($scope.models.dropzones[1][$rootScope.helpPopover.lastItem+1] == undefined) {
                        if ($rootScope.helpPopover.lastStep == "req_price") {
                            $rootScope.helpPopover.close()
                            return
                        } else {
                            $rootScope.helpPopover.lastStep = "req_price"
                            ind = $rootScope.helpPopover.lastItem.toString()+"req_price"
                            $rootScope.helpPopover.isOpen[ind] = true
                            $rootScope.helpPopover.content = $rootScope.help_texts[$rootScope.helpPopover.lastStep];
                            $rootScope.helpPopover.lastOpened = ind
                            return;
                        }
                    } else  {
                        $rootScope.helpPopover.close()
                        return
                    }
                } else {
                    if ($rootScope.helpPopover.lastStep == "dnd_button") {
                        ind = $rootScope.helpPopover.lastItem.toString()+
                            $rootScope.seq[$rootScope.helpPopover.lastStep];
                    } else if ($rootScope.helpPopover.lastStep == "inf_req") {
                        ind = $rootScope.helpPopover.lastItem.toString()+"|"+
                            $rootScope.helpPopover.lastItemPlus.toString()+
                            $rootScope.seq[$rootScope.helpPopover.lastStep];
                    } else {
                        ind = $rootScope.helpPopover.lastItem.toString()+"|"+
                            $rootScope.helpPopover.lastItemPlus.toString()+"|"+
                            $rootScope.helpPopover.lastMove.toString()+$rootScope.seq[$rootScope.helpPopover.lastStep];
                    }
                }
            } else {
                if ($rootScope.helpPopover.lastStep == "item") {
                    ind = '0dnd_button';
                } else {
                    ind = $rootScope.seq[$rootScope.helpPopover.lastStep]
                }
            }
            //define new step
            $rootScope.helpPopover.isOpen[ind] = $rootScope.helpPopover.showHelp;
            $rootScope.helpPopover.lastStep = $rootScope.seq[$rootScope.helpPopover.lastStep];
            $rootScope.helpPopover.content = $rootScope.help_texts[$rootScope.helpPopover.lastStep];
            $rootScope.helpPopover.lastOpened = ind
        },

        close: function close() {
            $rootScope.helpPopover.showHelp = false;
            $rootScope.helpPopover.lastStep = "";
            $rootScope.helpPopover.content = "";
            $rootScope.helpPopover.lastOpened = "";
            $rootScope.helpPopover.isOpen = {};
        }
    };
});

mTestApp.controller("mTestDBController",
    function ($scope, $timeout, $sce, $http, $cookies, $location, $routeParams, DB_data,
              $window, LS, mtCrud, MathData, ModalWin, $interval) {
        const token = localStorage.getItem('token');
        $scope.allowed = {dropzone: ['container'], container: ['itemplus', 'item'], itemplus: ['item']};

        $scope.params = $routeParams;

        var trusted = {};
        $scope.getPopoverContent = function (content) {
            return trusted[content] || (trusted[content] = $sce.trustAsHtml(content));
        };
        //----------------------------------------------------------------------------------------------------------
        //---------------------------------------------DATA---------------------------------------------------------
        //----------------------------------------------------------------------------------------------------------
        mtCrud.readMtestFromDB($scope.params.mtest_id)
            .then(function (response) {
                $scope.mtestData = response.data.mtest;
                $scope.executors = angular.fromJson(response.data.mtest.executors);
                this.val = angular.fromJson(response.data.mtest.calculations);
                $scope.modelsdb = {
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
                    dropzones: this.val
                };
            });

        //-------------------------------------------- FOR TYPEHEAD ------------------------------------------------

        $http({
            method: 'POST',
            url: '/api/v.1/m/regact/list',
            data: {mtest_id: $scope.params.mtest_id},
        }).then(function (response) {
            $scope.reg_files = response.data.reg_acts;
        });

        $http({
            method: 'GET',
            url: '/api/v.1/actions'
        }).then(function (response) {
            $scope.adm_actions = response.data.actions;
        });

        $http({
            method: 'GET',
            url: '/static/questions.json'
        }).then(function (response) {
            $scope.questions = response.data;
        }).catch(function (reason) {
            console.log(reason)
        });


        //-------------------------------------------- GROUP MATH ------------------------------------------------
        $scope.doGroupMath = function () {
            var execs = $scope.executors;
            var infs = $scope.modelsdb.dropzones['1'];

            for (var devs_get = []; devs_get.length < infs.length; devs_get.push([]));
            for (var p in execs) {
                let dev_data;
                console.log(execs[p].checked);
                if (execs[p].checked === true) {
                    $http({method: 'GET', url: '/api/v.1/m/get/' + execs[p].mid })
                        .then(function (response) {
                            dev_data = JSON.parse(response.data.mtest.calculations)
                        });
                    setTimeout(function () {
                        for (var j = 0; j < dev_data['1'].length; j++) {
                            if (dev_data['1'][j]['contsub'] !== 0) {
                                devs_get[j].push(dev_data['1'][j]['contsub'])
                            }
                        }
                    }, 50)
                }
            }
            $scope.devs_get = devs_get
        };

        $scope.calculateCorrTotal = function () {
            var total = 0;
            var infs = $scope.modelsdb.dropzones['1'];
            var arr = eval(infs);
            for(var i = 0; i < arr.length; i++) {
                var col = arr[i]['columns'][0];
                for(var j = 0; j < col.length; j++) {
                    if (col[j]['type'] === 'itemplus'){
                        var itemplus = col[j]['columns'][0];
                        for(var k = 0; k < itemplus.length; k++){
                            total += itemplus[k]['corr_calc']
                        }
                    }
                }
            }
            return total
        }

        setTimeout(function () {
            $scope.doGroupMath()
        }, 150);


        setTimeout(function () {
            $scope.sumDevsInfs = function (x) {
                var tmp_sum = 0;
                for (var i = 0; i < x.length; i++) {
                    tmp_sum += x[i]
                }

                if (isNaN(tmp_sum / x.length)) {
                    return 0
                }

                return tmp_sum / x.length
            }
        }, 200);


        //-------------------------------------------- END GROUP MATH ------------------------------------------------

        //-------------------------------corruption part

        $scope.corr_sum = function (data) {
            var total = 0;
            for (var k in data) {
                total += Number(data[k])
            }
            return total;
        };


        $scope.corr_bar = function (total, length) {
            if (total > length / 2) {
                $scope.barText = "Критична корупційна складова!";
                return 1;
            } else {
                $scope.barText = "Кількість корупційних ризиків:" + total;
                return 0;
            }
        };

        $scope.addRow = function (item, text) {
            var quest = {
                text: text,
                "yesno": [{"name": "Так", "value": "yes"}, {"name": "Ні", "value": "no"}],
                "remove": "block"
            };
            item.push(quest);
        };

        $scope.hideCritical = function (val, length) {
            if (val >= 100 || val > length / 2) {
                return "block"
            } else {
                return "none"
            }
        };

        $scope.removeRow = function (item, index) {
            item.splice(index, 1);
        };
        $scope.addToSlider = function (res, eft) {
            if (res == 1 && eft == 1) {
                return 1
            } else if (res == 1 && eft > 2) {
                return 100
            } else {
                return 0
            }
        };

        $scope.hideQuestion = function (dep, val) {
            if (dep == 1 && val == 0) {
                return "none"
            } else if (dep == 1 && val == 1) {
                return "block"
            } else {
                return "block"
            }
        };

        $scope.valueToText = function (val, type) {
            if (type == "yn") {
                if (val == 1 || val == "Yes") {
                    return "Так"
                } else if (val == 0 || val == "No") {
                    return "Ні"
                } else if (val == "idk") {
                    return "Незнаю"
                } else {
                    return "Не заповнено!"
                }
            } else {
                if (val == 1) {
                    return "Суттєво"
                } else if (val == 0) {
                    return "Несуттєво"
                } else if (val == 100) {
                    return "Критично"
                } else {
                    return "Не заповнено!"
                }
            }
        };

        //end corruption part --------------------

        //----------------------------------------------------------------------------------------------------------
        //---------------------------------------------SAVE TO DB --------------------------------------------------
        //----------------------------------------------------------------------------------------------------------

        $scope.mdata = function () {
            return angular.toJson($scope.modelsdb.dropzones);
        };
        $scope.saveMtestToDB = function () {
            var corr_total = $scope.calculateCorrTotal()
            mtCrud.updateMtestCalculations($scope.params.mtest_id, $scope.mdata(), corr_total, $scope.calc_total, token)
                .then(function (value) {
                    console.log(value.data)
                })
        };
        $scope.updateExecutors = function () {
            var executors = angular.toJson($scope.executors);
            mtCrud.updateMtestExecutors($scope.params.mtest_id, executors, token)
                .then(function (value) {
                    $scope.doGroupMath()
                    console.log(value.data)
                }).catch(function (err) {  console.log(err) })
        };
        //---------------------------------------------SUM----------------------------------------------------------

        $scope.SumDb = function (id) {
            return MathData.mathSum($scope.modelsdb.dropzones['1'], id)
        };
        $scope.totalSumDb = function () {
            $scope.calc_total = MathData.totalSum($scope.modelsdb.dropzones['1'][0]['ki'], $scope.SumDb)
            return $scope.calc_total
        };

        $scope.awgSumDb = function () {
            $scope.calc_total = MathData.awgSum($scope.modelsdb.dropzones['1'], $scope.modelsdb.dropzones['1'][0]['ki'])
            return MathData.awgSum($scope.modelsdb.dropzones['1'], $scope.modelsdb.dropzones['1'][0]['ki'])
        };

        //---------------------------------------------MODAL--------------------------------------------------------
        $scope.openModal = function (size, template, controller) {
            return ModalWin.openModal(size, template, controller)
        };


        //---------------------------------------------alert------------------------------------------------------


        $scope.getClass = 'not_disabled';
        $scope.getClassTwo = 'disabled';

        $scope.setClasses = function () {
            $scope.getClass = 'disabled';
            $scope.getClassTwo = 'not_disabled';
            $timeout(function () {
                $scope.getClass = 'not_disabled';
                $scope.getClassTwo = 'disabled';
            }, 1000)
        };

        $scope.getRegAct = function (mtestID, docID) {
            console.log(mtestID, docID);
            $http({
                method: 'POST',
                url: "/api/v.1/m/regact/get",
                data: {mtest_id: mtestID, doc_id: docID},
                headers: {
                    'Content-Type': 'application/json', Authorization: 'Bearer ' + token
                }
            }).then(function (response) {
                console.log(response)
            }).catch(function (err) {
                console.log(err)
            });
        };

        //---------------------------------------------PDF--------------------------------------------------------
        $interval(function () {
            var a = document.getElementById('report').innerHTML;
            var blob = new Blob(["\ufeff", a], {
                type: 'text/html'
            });
            $scope.saveToPdf = (window.URL || window.webkitURL).createObjectURL(blob);
        }, 1000);


//---------------------------------------------reset------------------------------------------------------
        $scope.resetdb = function () {
            mtCrud.readMtestFromDB($scope.params.mtest_id).then(function (response) {
                var val = angular.fromJson(response.data.mtest.calculations);
                $scope.modelsdb.dropzones = angular.copy(
                    val
                )
            })
        };

        $scope.effectPopover = $sce.trustAsHtml('Критично – реалізація корупційного ризику призведе до зупинки діяльності <br> ' +
            'Суттєво - реалізація корупційного ризику призведе до суттєвих витрат грошей та/або часу співставних з розміром щомісячного прибутку');
        $scope.payContent = $sce.trustAsHtml('Оплата – зазвичай середня  вартість людино-години роботи персоналу, який\
     виконує дії необхідні для  виконання ІВ.<br>\
  За бажанням можна використовувати  заробітну плату однієї години відповідного Спеціаліста або Керівника. <br>\
  Для ФОПів  це може бути  розрахунок вартості  людино-години на основі річного обсягу реалізації сектору (детальніше див.  Посібник з М-Тесту).<br>\
  В віконце  оплата вводиться  за місяць і перераховується автоматично в оплату за годину.');
    });

mTestApp.controller("authLoginController", function ($scope, $timeout, $location, authService) {
    $scope.template_name = true;
    $scope.user = {};
    $scope.onLogin = function () {
        authService.login($scope.user)
            .then(function (user) {
                localStorage.setItem('token', user.data.token);
                $location.path('/u/cabinet');
            })
            .catch(function (err) {
                console.log(err);
                $scope.message = "Невірний логін або пароль, спробуйте ще раз";
                $timeout(function () {
                    $scope.message = ""
                }, 2000);
            });
    };
    $scope.sendResetRequest = function () {
        authService.resetpass($scope.resetpass.email)
            .then(function () {
                $scope.message = "На ваш Email відправлено посилання для відновлення паролю"
            })
            .catch(function () {
                $scope.message = "Даний Email не зареєстровано";
                $timeout(function () {
                    $scope.message = ""
                }, 2000);
            })
    }
});

mTestApp.controller("userCabinetController", function ($scope, $http, $location, $rootScope,
                                                       $timeout, authService, mtCrud, ModalWin, fileUploadService) {
    $scope.changepass = false;
    const token = localStorage.getItem('token');
    if (token) {
        authService.ensureAuthenticated(token, 'api/v.1/u/cabinet')
            .then(function (user) {
                if (user.status === 200) {
                    $scope.userdata = user.data.data;
                    $scope.records = user.data.data.records;
                    $rootScope.isLoggedIn = true;
                }
            })
            .catch(function (err) {
                console.log(err);
                $location.path('/u/login');
            });
    }



    //load governments and regions
    $http({
        method: 'GET',
        url: '/api/v.1/regions',
    }).then(function (response) {
        $scope.regions = response.data.regions
    }).catch( function (reason) {
        console.log(reason)
        });

    $http({
        method: 'GET',
        url: '/api/v.1/governments',
    }).then(function (response) {
        $scope.governs = response.data.govs
    }).catch( function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/businesses',
    }).then(function (response) {
        $scope.businesses = response.data.businesses
    }).catch(function (reason) {
        console.log(reason)
    });

    //end load governments and regions

    //format label for typehead on select
    $scope.formatLabel = function(model, index, itmtype) {
        console.log(model, index);
        for (var i=0; i< $scope[itmtype+'s'].length; i++) {
            if (model.id === $scope[itmtype+'s'][i].id) {
                $scope.records[index][itmtype] = $scope[itmtype+'s'][i].id
            }
        }
    };
    //end format label for typehead on select

    $scope.changeUserField = function (field, id, value) {
        console.log(field, id, value);
        $http({
            method: 'POST',
            url: "/api/v.1/u/edituser",
            data: {field: field, data: value, id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.openModal = function (size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };
    //--------------CRUD mtest------------------------
    //add executor to mtest item
    $scope.addExecutor = function (title, email, region, gov, dev_mid) {
        mtCrud.addMtestExecutor(title, email,region, gov, dev_mid, token).then(function (response) {
            $scope.records[dev_mid].executors[email] = {email: email, mid: response.data.mid}
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeExecutor = function (email, exIndex, devIndex) {
        mtCrud.removeMtestExecutor(email, exIndex, devIndex, token).then(function () {
            delete $scope.records[devIndex].executors[exIndex]
        }).catch(function (err) {
            console.log(err)
        });
    };
    //end add executor to mtest item
    $scope.addMtest = function (newmtest) {
        mtCrud.addMtest(newmtest, token).then(function (response) {
            $scope.records = response.data.records;
            $scope.len_records = angular.toJson($scope.records).length
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeMtestItem = function (id) {
        console.log(id);
        mtCrud.removeMtestItem(id, token)
            .then(function () {
                delete $scope.records[id];
                $scope.len_records = angular.toJson($scope.records).length
            }).catch(function (err) {
            console.log(err)
        });
    };

    $scope.updateMtestItem = function (item) {
        mtCrud.updateMtestItem(item, token)
            .then(function (response) {
            }).catch(function (err) {
            console.log(err)
        });
    };

    $scope.uploadRegAct = function (images, mtestID) {
        var uploadUrl = "/api/v.1/m/regact",
            promise = fileUploadService.uploadFileToUrl(images, uploadUrl, mtestID, token);

        promise.then(function (response) {
            if ($scope.records[mtestID].files == null) {
                $scope.records[mtestID].files = []
            }
            $scope.records[mtestID].files.push({
                "DocID": response.data.act.doc_id,
                "MtestID": mtestID,
                "Name": response.data.act.name,
                "Type": response.data.act.type
            })
        }, function () {
            $scope.serverResponse = 'An error has occurred';
        })
    };
    $scope.getRegAct = function (mtestID, docID) {
        console.log(mtestID, docID);
        $http({
            method: 'POST',
            url: "/api/v.1/m/regact/get",
            data: {mtest_id: mtestID, doc_id: docID},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            console.log(response)
        }).catch(function (err) {
            console.log(err)
        });
    };

    $scope.removeRegAct = function (mtestID, docID) {
        console.log(mtestID, docID);
        const index = $scope.records[mtestID].files.findIndex(a => a.docID === docID)
        $http({
            method: 'DELETE',
            url: "/api/v.1/m/regact",
            data: {mtest_id: mtestID, doc_id: docID},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.records[mtestID].files.splice(index,1)
        }).catch(function (err) {
            console.log(err)
        });
    };

    // popover for setting, mail, add mtest etc..
    $scope.dynamicPopover = {
        isOpen: {},
        title: "heh",

        templateUrl: {},

        open: function open(index, template, tp) {
            $scope.dynamicPopover.templateUrl[index + tp] = template;
            $scope.dynamicPopover.isOpen[index + tp] = true;
            $scope.dynamicPopover.data = 'Hello!';
        },

        close: function close(index, tp) {
            $scope.dynamicPopover.isOpen[index + tp] = false;
        }
    };
    // end popover for setting, mail, add mtest etc..

    //count len of records array for showing-hiding empty array banner
    var countUp = function () {
        $scope.len_records = angular.toJson($scope.records).length
    };
    $timeout(countUp, 300);
    //end count len of records array for showing-hiding empty array banner
});

mTestApp.controller("authRegisterController", function ($scope, authService, $location,$timeout) {
    $scope.user = {
        password: "",
        confirmPassword: ""
    };
    $scope.onRegister = function () {
        authService.register($scope.user)
            .then(function (response) {
                $location.path('/u/login');
            })
            .catch(function (err) {
                $scope.show_err = 1;
                $scope.err_status = err.data.title;
                $timeout(function () {
                    $scope.show_err = 0
                }, 2000);
                console.log(err);
            });
    };
});

mTestApp.controller("menuController", function ($scope, $rootScope, $location, authService, ModalWin) {
    $rootScope.isLoggedIn = false;
    const token = localStorage.getItem('token');
    if (token) {
        authService.ensureAuthenticated(token, 'api/v.1/u/cabinet')
            .then(function (user) {
                if (user.status === 200) {
                    $rootScope.isLoggedIn = true;
                }
            })
            .catch(function () {
                $rootScope.isLoggedIn = false;
            });
    }
    $scope.onLogout = function () {
        localStorage.removeItem('token');
        $rootScope.isLoggedIn = false;
        $location.path('/u/login');
    };
    $scope.openModal = function (size, template, controller) {
        return ModalWin.openModal(size, template, controller)
    };
});

mTestApp.controller("searchController", function ($scope, $http) {
    $scope.currentPage = 0;
    //load governments and regions
    $http({
        method: 'GET',
        url: '/api/v.1/regions',
    }).then(function (response) {
        $scope.regions = response.data.regions
    }).catch( function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/governments',
    }).then(function (response) {
        $scope.governs = response.data.govs
    }).catch( function (reason) {
        console.log(reason)
    });
    $http({
        method: 'GET',
        url: '/api/v.1/businesses',
    }).then(function (response) {
        $scope.businesses = response.data.businesses
    }).catch( function (reason) {
        console.log(reason)
    });
    //end load governments and regions
    //elastic search
    $scope.query = {
        "from": 0, "size": 10,
        "query": {
            "bool": {
                "should": {
                    "multi_match": {
                        "query": $scope.phrase,
                        "fields": ["name", "reg_act"]
                    }
                }
            }
        }
    };

    $scope.addPhrase = function () {
        $scope.query.query.bool.should.multi_match.query = $scope.phrase
    };

    $scope.addTerm = function (field, data) {
        if (!$scope.query.query.bool.filter) {
            $scope.query.query.bool.filter = {
                "bool": {
                    "must": []
                }
            }
        }
        var obj = {};
        var arr = $scope.query.query.bool.filter.bool.must;
        obj[field] = data;
        if (arr.length === 0) {
            arr.push({"term": obj})
        } else {
            for (var i = 0; i < arr.length; i++) {
                if (arr[i].term.hasOwnProperty(field)) {
                    arr[i].term[field] = data;
                    break
                } else if (i === arr.length - 1 && !arr[i].term.hasOwnProperty(field)) {
                    arr.push({"term": obj})
                }
            }
        }
    };

    $scope.propertyName = 'corr_result';
    $scope.reverse = true;

    $scope.sortBy = function(propertyName) {
        $scope.reverse = ($scope.propertyName === propertyName) ? !$scope.reverse : false;
        $scope.propertyName = propertyName;
    };

    $scope.doSearch = function () {
        $http({
            method: 'POST',
            url: "http://mtest.org.ua/api/v.1/search",
            data: $scope.query
        }).then(function (response) {
            $scope.results = response.data;
        }).catch(function (err) {
            console.log(err)
        })
    };

    $scope.getBusinessNames = function (id) {
        for (var i = 0; i < $scope.businesses.length; i++) {
            if ($scope.businesses[i].id == id) {
                return $scope.businesses[i].name
            }
        }
        return "Назву бізнесу не вказано"
    }
});

mTestApp.controller("authActivateController", function ($scope, $routeParams,$http) {
    const baseURL = 'http://mtest.org.ua';
    $http({
        method: 'GET',
        url: baseURL + 'api/v.1/u/activate/' + $routeParams.hash ,
    }).then(function (response) {
        $scope.message = "Акаунт успішно активовано"
    }).catch(function (err) {
        $scope.err_message = "Посилання не існує"
    });

});

mTestApp.controller("authResetController", function ($scope, $routeParams, $http, $location, authService) {
    const baseURL = 'http://mtest.org.ua';
    var hash = $routeParams.hash;
    $scope.user = {};
    authService.checkhash(hash)
        .then(function (value) {
        })
        .catch(function (reason) {
            console.log(reason);
            $location.url('/');
        });
    $scope.onReset = function () {
        authService.setnewpass($scope.user.password, hash)
            .then(function () {
                $scope.message = "Пароль успішно змінено"
            }).catch(function () {
            $scope.err_message = "Помилка: посилання не існує"
        });
    }

});


mTestApp.controller("adminController", function ($scope, $http, $rootScope ,$location, authService) {
    $scope.changepass = false;
    const token = localStorage.getItem('token');
    if (token) {
        authService.ensureAuthenticated(token, 'api/v.1/admin')
            .then(function (user) {
                if (user.status === 200) {
                    $rootScope.isLoggedIn = true;
                }
            })
            .catch(function (err) {
                console.log(err);
                $location.path('/u/login');
            });
    }
    //load governments and regions
    $http({
        method: 'GET',
        url: '/api/v.1/regions',
    }).then(function (response) {
        $scope.regions = response.data.regions
    }).catch( function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/businesses',
    }).then(function (response) {
        $scope.businesses = response.data.businesses
    }).catch(function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/synonyms',
    }).then(function (response) {
        $scope.synonyms = response.data.synonyms
        console.log($scope.synonyms)
    }).catch(function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/governments',
    }).then(function (response) {
        $scope.governments = response.data.govs
    }).catch( function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/actions',
    }).then(function (response) {
        $scope.actions = response.data.actions
    }).catch( function (reason) {
        console.log(reason)
    });

    $http({
        method: 'GET',
        url: '/api/v.1/users',
    }).then(function (response) {
        $scope.users = response.data.users
    }).catch( function (reason) {
        console.log(reason)
    });
    // popover for setting, mail, add mtest etc..
    $scope.dynamicPopover = {
        isOpen: {},
        title: "heh",

        templateUrl: {},

        open: function open(index, template, tp) {
            $scope.dynamicPopover.templateUrl[index + tp] = template;
            $scope.dynamicPopover.isOpen[index + tp] = true;
            $scope.dynamicPopover.data = 'Hello!';
        },

        close: function close(index, tp) {
            $scope.dynamicPopover.isOpen[index + tp] = false;
        }
    };
    $scope.changeAdminUserField = function (field, id, value) {
        console.log(field, id, value);
        $http({
            method: 'POST',
            url: "/api/v.1/u/edituser",
            data: {field: field, data: value, id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeUser = function (id) {
        console.log(id);
        const index =  $scope.users.findIndex(a => a.id === parseInt(id))
        $http({
            method: 'DELETE',
            url: "/api/v.1/u/delete",
            data: { id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.users.splice(index,1)
        }).catch(function (err) {
            console.log(err)
        });
    };
    //region
    $scope.saveRegionName = function (id, value) {
        console.log(id, value);
        $http({
            method: 'PUT',
            url: "/api/v.1/regions",
            data: {name: value, id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function (err) {
            console.log(err)
        });
    };
    //governments
    $scope.saveGovernmentName = function (id, value) {
        console.log(id, value);
        $http({
            method: 'PUT',
            url: "/api/v.1/governments",
            data: {name: value, id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.addGovernment = function (value) {
        console.log(value);
        $http({
            method: 'POST',
            url: "/api/v.1/governments",
            data: {name: value},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.governments.push({name:value, id: $scope.governments.length+1})
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeGovernment = function (value) {
        console.log(value);
        const index =  $scope.governments.findIndex(a => a.id === parseInt(value))
        $http({
            method: 'DELETE',
            url: "/api/v.1/governments",
            data: {id: parseInt(value)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.governments.splice(index,1)
        }).catch(function (err) {
            console.log(err)
        });
    };
    ///actions
    $scope.saveActionName = function (id, value) {
        console.log(id, value);
        $http({
            method: 'PUT',
            url: "/api/v.1/actions",
            data: {name: value, id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.addAction = function (value) {
        console.log(value);
        $http({
            method: 'POST',
            url: "/api/v.1/actions",
            data: {name: value},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.actions.push({act_name:value, id: $scope.actions.length+1})
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeAction = function (value) {
        console.log(value);
        const index =  $scope.actions.findIndex(a => a.id === parseInt(value))
        $http({
            method: 'DELETE',
            url: "/api/v.1/actions",
            data: {id: parseInt(value)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.actions.splice(index,1)
        }).catch(function (err) {
            console.log(err)
        });
    };

    ///actions
    $scope.saveBusinessName = function (id, value) {
        console.log(id, value);
        $http({
            method: 'PUT',
            url: "/api/v.1/businesses",
            data: {name: value, id: parseInt(id)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.addBusinessType = function (value) {
        console.log(value);
        $http({
            method: 'POST',
            url: "/api/v.1/businesses",
            data: {name: value},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.businesses.push({name:value, id: $scope.businesses.length+1})
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeBusinessType = function (value) {
        console.log(value);
        const index =  $scope.businesses.findIndex(a => a.id === parseInt(value))
        $http({
            method: 'DELETE',
            url: "/api/v.1/businesses",
            data: {id: parseInt(value)},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.businesses.splice(index,1)
        }).catch(function (err) {
            console.log(err)
        });
    };

    $scope.addSynonym = function (word, synonym) {
        console.log(word, synonym);
        $http({
            method: 'POST',
            url: "/api/v.1/synonyms",
            data: {word: word, synonym:synonym},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.synonyms.push({word:word, synonym: synonym})
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeSynonym = function (word, synonym) {
        console.log(word, synonym);
        const index =  $scope.businesses.findIndex(a => a.word === word && a.synonym === synonym)
        $http({
            method: 'DELETE',
            url: "/api/v.1/synonyms",
            data: {word: word, synonym:synonym},
            headers: {
                'Content-Type': 'application/json', Authorization: 'Bearer ' + token
            }
        }).then(function (response) {
            $scope.businesses.splice(index,1)
        }).catch(function (err) {
            console.log(err)
        });
    };
    $scope.removeUserMtest = function (id) {
        console.log(id);
        mtCrud.removeMtestItem(id, token)
            .then(function () {
                delete $scope.records[id];
                $scope.len_records = angular.toJson($scope.records).length
            }).catch(function (err) {
            console.log(err)
        });
    };

});
