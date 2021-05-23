var mTestApp = angular.module('mTestApp');
mTestApp.factory('ModalWin', function ($uibModal, $log) {
    return {
        openModal: function (size,template,controller) {
        var modalInstance = $uibModal.open({
          animation: true,
          ariaLabelledBy: 'modal-title',
          ariaDescribedBy: 'modal-body',
          templateUrl: template,
          controller: controller,
          controllerAs: 'm',
          size: size
        });
        modalInstance.result.then(function (selectedItem) {
          m.selected = selectedItem;
        }, function () {
          $log.info('Modal dismissed at: ' + new Date());
        })
      }
       
    }
});

