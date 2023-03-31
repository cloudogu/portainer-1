import { useQuery } from 'react-query';

import { getBackupStatus } from '@/portainer/services/api/backup.service';
import { isoDate } from '@/portainer/filters/filters';

import { InformationPanel } from '@@/InformationPanel';
import { TextTip } from '@@/Tip/TextTip';
import { Link } from '@@/Link';

export function BackupFailedPanel() {
  const { status, isLoading } = useBackupStatus();

  if (isLoading || !status || !status.Failed) {
    return null;
  }

  return (
    <InformationPanel title="Information">
      <TextTip>
        The latest automated backup has failed at {isoDate(status.TimestampUTC)}
        . For details please see the log files and have a look at the{' '}
        <Link to="portainer.settings">settings</Link> to verify the backup
        configuration.
      </TextTip>
    </InformationPanel>
  );
}

function useBackupStatus() {
  const { data, isLoading } = useQuery(
    ['backup', 'status'],
    () => getBackupStatus(),
    {
      onError() {
        // ignore license notifications
      },
    }
  );

  return { status: data, isLoading };
}
