import React from 'react';
import PropTypes from 'prop-types';

import GroupHeader from './GroupHeader';
import OutletList from './OutletList';

export default function GroupListItem(props) {
  const { id, name, outlets, dispatchMessage } = props;

  const handleAction = (action) => () => {
    dispatchMessage({ type: 'group', data: { id, action } });
  }

  return (
    <React.Fragment>
      <GroupHeader
        name={name}
        onActionOn={handleAction('on')}
        onActionOff={handleAction('off')}
        onActionToggle={handleAction('toggle')}
      />
      <OutletList outlets={outlets} dispatchMessage={dispatchMessage} />
    </React.Fragment>
  );
}

GroupListItem.propTypes = {
 id: PropTypes.string.isRequired,
 name: PropTypes.string.isRequired,
 outlets: PropTypes.array.isRequired,
 dispatchMessage: PropTypes.func.isRequired,
};
