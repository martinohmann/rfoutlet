import React from 'react';
import PropTypes from 'prop-types';

import GroupHeader from './GroupHeader';
import OutletList from './OutletList';
import websocket from '../websocket';

export default function GroupListItem(props) {
  const { id, name, outlets } = props;

  const handleAction = (action) => () => {
    websocket.sendMessage({ type: 'group', data: { id, action } });
  }

  return (
    <React.Fragment>
      <GroupHeader
        name={name}
        onActionOn={handleAction('on')}
        onActionOff={handleAction('off')}
        onActionToggle={handleAction('toggle')}
      />
      <OutletList outlets={outlets} />
    </React.Fragment>
  );
}

GroupListItem.propTypes = {
 id: PropTypes.string.isRequired,
 name: PropTypes.string.isRequired,
 outlets: PropTypes.array.isRequired,
};
