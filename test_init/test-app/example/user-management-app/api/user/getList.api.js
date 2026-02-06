/**
 * @api
 * @name getList
 * @path /api/user/getList
 * @method POST
 * @description Get user list
 * @group user
 * @version 1.0.0
 */

// @param
// name: pageNum
// type: Integer
// required: true
// default: 1
// description: Page number

// @param
// name: pageSize
// type: Integer
// required: true
// default: 10
// description: Page size

// @return
// type: PageResult
// description: Paginated user list

(function() {
    var pageNum = parseInt($params.pageNum || 1);
    var pageSize = parseInt($params.pageSize || 10);

    var countResult = $db.query("SELECT COUNT(*) as total FROM platform_user");
    var total = countResult[0].total;

    var offset = (pageNum - 1) * pageSize;
    var list = $db.query(
        "SELECT id, name, login_name, email, status, created_at FROM platform_user ORDER BY created_at DESC LIMIT ? OFFSET ?",
        [pageSize, offset]
    );

    return {
        code: 200,
        message: "success",
        data: {
            list: list,
            total: total,
            pageNum: pageNum,
            pageSize: pageSize,
            pages: Math.ceil(total / pageSize)
        }
    };
})();
