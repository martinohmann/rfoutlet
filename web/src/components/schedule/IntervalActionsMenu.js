import React from 'react';
import PropTypes from 'prop-types';
import IconButton from '@material-ui/core/IconButton';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import Typography from '@material-ui/core/Typography';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import EditIcon from '@material-ui/icons/Edit';
import DeleteIcon from '@material-ui/icons/Delete';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import { useTranslation } from 'react-i18next';

export default function IntervalActionsMenu({ onDelete, onEdit }) {
  const [anchorElement, setAnchorElement] = React.useState();

  const handleEdit = () => {
    setAnchorElement(null);
    onEdit();
  }

  const handleDelete = () => {
    setAnchorElement(null);
    onDelete();
  }

  const { t } = useTranslation();

  return (
    <span>
      <IconButton
        aria-owns={anchorElement ? 'interval-actions-menu' : undefined}
        aria-haspopup="true"
        onClick={(e) => setAnchorElement(e.currentTarget)}
      >
        <MoreVertIcon />
      </IconButton>
      <Menu
        id="interval-actions-menu"
        anchorEl={anchorElement}
        open={Boolean(anchorElement)}
        onClose={() => setAnchorElement(null)}
      >
        <MenuItem onClick={handleEdit}>
          <ListItemIcon>
            <EditIcon />
          </ListItemIcon>
          <Typography variant="inherit" noWrap>
            {t('edit')}
          </Typography>
        </MenuItem>
        <MenuItem onClick={handleDelete}>
          <ListItemIcon>
            <DeleteIcon />
          </ListItemIcon>
          <Typography variant="inherit" noWrap>
            {t('delete')}
          </Typography>
        </MenuItem>
      </Menu>
    </span>
  );
}

IntervalActionsMenu.propTypes = {
  onDelete: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
};
