export function errorForStatus(res: any, catchErr: string): Error {
  if (res.status === 401) {
    const err = new Error('Unauthorised access to API')
    err.name = 'AuthError'
    return err
  } else if (res.status === 403) {
    const err = new Error('You do not have permission to access resource')
    return err
  } else if (res.status === 404) {
    return new Error('Resource not found')
  } else if (res.status === 500) {
    return new Error('Internal server error')
  } else {
    return new Error(catchErr)
  }
}
