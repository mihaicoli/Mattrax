export function errorForStatus(res: any, catchErr: string): Error {
  if (res.status === 401) {
    return new Error('Unauthorised access to API')
  } else if (res.status === 404) {
    return new Error('Resource not found')
  } else if (res.status === 500) {
    return new Error('Internal server error')
  } else {
    return new Error(catchErr)
  }
}
