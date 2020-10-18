// TODO: Pagination, Filters

export const actions = {
  getAll(_context: any) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/a1663929-dfc6-4801-8052-769fc15eb79d')
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
  getByID(_context: any, _userID: string) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/d779e60b-f5d8-4931-bc4c-03b09d16d587')
        .then(async (res) => {
          if (res.status !== 200) {
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
