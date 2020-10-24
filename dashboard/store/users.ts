// TODO: Pagination, Filters

export const actions = {
  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/users', {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching users from server'))
            return
          }

          const users = await res.json()
          resolve(users)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(context: any, userID: string) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/user/' + encodeURI(userID), {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status === 404) {
            resolve(null)
            return
          } else if (res.status !== 200) {
            reject(new Error('Error fetching user from server'))
            return
          }

          const user = await res.json()
          resolve(user)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
