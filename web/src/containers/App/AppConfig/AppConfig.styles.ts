import { css } from '@emotion/react'

import { ThemeType, get } from '@shared-ui/components/Atomic/_theme'

export const item = (theme: ThemeType) => css`
    display: flex;
    align-items: center;
    color: ${get(theme, `AppConfig.icon.color`)};
    transition: all 0.3s;
    cursor: pointer;

    &:hover {
        color: ${get(theme, `AppConfig.icon.hover.color`)};
    }
`

export const floatingPanel = (theme: ThemeType) => css`
    z-index: 10;
    background: ${get(theme, `NotificationCenter.floatingPanel.background`)};
    border: 1px solid ${get(theme, `NotificationCenter.floatingPanel.borderColor`)};
    box-shadow: 0 30px 40px rgba(28, 52, 99, 0.1);
    border-radius: 8px;
    min-width: 400px;
    max-width: 600px;
`
