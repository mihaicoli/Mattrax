// TODO: Pagination, Filters

export const actions = {
  getAll(_context: any) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/6d4001df-2d95-4cf0-82ff-3271f27e2dd5')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching policies from server'))
            return
          }

          const policies = await res.json()
          resolve(policies)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(_context: any, _policyID: string) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/771ff025-55c6-4845-a255-9a095333af3a')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching policy from server'))
            return
          }

          const policy = await res.json()
          resolve(policy)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
