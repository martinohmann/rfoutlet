import React from 'react';
import PropTypes from 'prop-types';

import GroupListItem from './GroupListItem';
import { List } from './List';

export default function GroupList(props) {
  return (
    <List>
      {props.groups.map(group =>
        <GroupListItem key={group.id} {...group} />
      )}
    </List>
  );
}

GroupList.propTypes = {
 groups: PropTypes.array.isRequired,
};
