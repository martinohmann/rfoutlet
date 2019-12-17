import React from 'react';
import PropTypes from 'prop-types';
import { List, NoItemsListItem } from '../List';
import { useTranslation } from 'react-i18next';
import GroupHeader from './GroupHeader';
import OutletList from './OutletList';
import dispatcher from '../../dispatcher';

export default function GroupList({ groups }) {
  const { t } = useTranslation();

  return (
    <List>
      {groups.map(group =>
        <GroupListItem key={group.id} {...group} />
      )}
      {groups.length === 0 && (
        <NoItemsListItem
          primary={t('no-groups-primary')}
          secondary={t('no-groups-secondary')}
        />
      )}
    </List>
  );
}

GroupList.propTypes = {
  groups: PropTypes.array.isRequired,
};

const GroupListItem = ({ id, name, outlets }) => {
  const handleAction = (action) => () => {
    dispatcher.dispatchGroupMessage(id, action);
  };

  return (
    <>
      <GroupHeader
        name={name}
        onActionOn={handleAction('on')}
        onActionOff={handleAction('off')}
        onActionToggle={handleAction('toggle')}
      />
      <OutletList outlets={outlets} />
    </>
  );
}

GroupListItem.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  outlets: PropTypes.array.isRequired,
};
