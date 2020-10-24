export default function (context: any) {
  context.store
    .dispatch('authentication/populateUserInfomation')
    .then(() => context.store.dispatch('authentication/isAuthenticated'))
    .then((authenticated: boolean) => {
      if (!authenticated) {
        context.app.router.push('/login')
      }
    })
}
