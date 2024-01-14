export const getUser = store =>
  store && store.user ? store.user : {};

export const getBots = (store, id) =>
  store && store.bots
    ? store.bots
    : {};