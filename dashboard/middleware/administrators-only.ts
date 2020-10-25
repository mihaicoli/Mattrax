export default function (context: any) {
  if (context.store.state.authentication.user.aud === 'enrollment') {
    context.app.router.push('/enroll')
  }
}
