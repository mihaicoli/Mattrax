// TODO: Pagination, Filters

export const actions = {
  getAll(_context: any) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/0a624a65-eb2c-4ca5-aeb2-822e6397dfd0')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching devices from server'))
            return
          }

          const devices = await res.json()
          resolve(devices)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(_context: any, _deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch('https://run.mocky.io/v3/45a4d87a-92fc-4f30-9ca3-5ef2bf284385')
        .then(async (res) => {
          if (res.status !== 200) {
            reject(new Error('Error fetching device from server'))
            return
          }

          const device = await res.json()
          resolve(device)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
