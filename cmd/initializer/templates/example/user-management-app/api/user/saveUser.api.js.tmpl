/**
 * @api
 * @name saveUser
 * @path /api/user/saveUser
 * @method POST
 * @description Save user
 * @group user
 * @version 1.0.0
 */

// @param
// name: id
// type: Integer
// required: false
// description: User ID (empty for new user)

// @param
// name: name
// type: String
// required: true
// description: User name

(function() {
    var id = $params.id;
    var name = $params.name;
    var loginName = $params.loginName;

    if (!name || !loginName) {
        return { code: 400, message: "Name and login name are required", data: null };
    }

    if (id) {
        $db.execute("UPDATE platform_user SET name = ?, login_name = ? WHERE id = ?", [name, loginName, parseInt(id)]);
    } else {
        $db.execute("INSERT INTO platform_user (name, login_name) VALUES (?, ?)", [name, loginName]);
    }

    return { code: 200, message: "success", data: { success: true } };
})();
