import { useEffect } from 'react'
import PropTypes from 'prop-types'
import BPagination from 'react-bootstrap/Pagination'
import { useIntl } from 'react-intl'
import classNames from 'classnames'

import { PaginationItems } from './pagination-items'
import { messages as t } from './table-i18n'

export const Pagination = props => {
  const {
    className,
    disabled,
    canPreviousPage,
    canNextPage,
    pageCount,
    gotoPage,
    nextPage,
    previousPage,
    pageIndex,
    pageLength,
  } = props

  const { formatMessage: _ } = useIntl()

  // If the last item is removed from the list, and we are on the last page (pageLength === 0), update the last page with (pageCount - 1)
  // Only do this if there are the least 2 pages available (pageCount > 1)
  useEffect(() => {
    if (pageLength === 0 && pageCount >= 1) {
      const prevPage = pageCount - 1
      gotoPage(prevPage >= 0 ? prevPage : 0)
    }
  }, [gotoPage, pageCount, pageLength])

  return (
    <BPagination
      className={classNames('plgd-pagination', className)}
      disabled={disabled}
    >
      {/* <BPagination.First onClick={() => gotoPage(0)} disabled={!canPreviousPage} /> */}
      <BPagination.Prev
        className="step"
        onClick={() => previousPage()}
        disabled={!canPreviousPage}
      >
        {_(t.prev)}
      </BPagination.Prev>
      <PaginationItems
        activePage={pageIndex + 1}
        pageCount={pageCount}
        maxButtons={10}
        onItemClick={gotoPage}
      />
      <BPagination.Next
        className="step"
        onClick={() => nextPage()}
        disabled={!canNextPage}
      >
        {_(t.next)}
      </BPagination.Next>
      {/* <BPagination.Last onClick={() => gotoPage(pageCount - 1)} disabled={!canNextPage} /> */}
    </BPagination>
  )
}

Pagination.propTypes = {
  className: PropTypes.string,
  disabled: PropTypes.bool,
  canPreviousPage: PropTypes.bool.isRequired,
  canNextPage: PropTypes.bool.isRequired,
  pageCount: PropTypes.number.isRequired,
  gotoPage: PropTypes.func.isRequired,
  nextPage: PropTypes.func.isRequired,
  previousPage: PropTypes.func.isRequired,
  setPageSize: PropTypes.func.isRequired,
  pageIndex: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  pageSizes: PropTypes.arrayOf(PropTypes.number),
  pageLength: PropTypes.number.isRequired,
}

Pagination.defaultProps = {
  className: null,
  disabled: false,
  pageSizes: [10, 20, 30, 40, 50],
}
