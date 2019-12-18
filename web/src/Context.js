import React from 'react';

const Context = React.createContext({});

export default Context;

export function GroupProvider({ groups, children, ...rest }) {
  const outlets = groups.reduce((outlets, group) => {
    return outlets.concat(group.outlets);
  }, []);

  const intervals = outlets.reduce((intervals, outlet) => {
    return intervals.concat(outlet.schedule);
  }, []);

  const value = { groups, outlets, intervals };

  return (
    <Context.Provider value={value} {...rest}>
      {children}
    </Context.Provider>
  );
}
