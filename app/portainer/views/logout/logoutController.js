import angular from 'angular';

class LogoutController {
  /* @ngInject */
  constructor($async, $state, $transition$, Authentication, StateManager, Notifications, LocalStorage, SettingsService) {
    this.$async = $async;
    this.$state = $state;
    this.$transition$ = $transition$;

    this.Authentication = Authentication;
    this.StateManager = StateManager;
    this.Notifications = Notifications;
    this.LocalStorage = LocalStorage;
    this.SettingsService = SettingsService;

    this.logo = this.StateManager.getState().application.logo;
    this.logoutAsync = this.logoutAsync.bind(this);

    this.onInit = this.onInit.bind(this);
  }

  /**
   * UTILS FUNCTIONS SECTION
   */
  async logoutAsync() {
    const error = this.$transition$.params().error;
    const performApiLogout = this.$transition$.params().performApiLogout;
    try {
      await this.Authentication.logout(performApiLogout);
      const settings = await this.SettingsService.publicSettings();
      document.location.href = settings.OAuthLogoutURI;
    } finally {
      this.LocalStorage.storeLogoutReason(error);
    }
  }

  logout() {
    return this.$async(this.logoutAsync);
  }

  /**
   * END UTILS FUNCTIONS SECTION
   */

  async onInit() {
    try {
      await this.logout();
    } catch (err) {
      this.Notifications.error('Failure', err, 'An error occured during logout');
    }
  }

  $onInit() {
    return this.$async(this.onInit);
  }

  /**
   * END ON INIT SECTION
   */
}

export default LogoutController;
angular.module('portainer.app').controller('LogoutController', LogoutController);
