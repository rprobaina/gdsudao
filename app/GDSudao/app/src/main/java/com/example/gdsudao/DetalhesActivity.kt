package com.example.gdsudao

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.widget.Toast

class DetalhesActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_detalhes)

        setSupportActionBar(findViewById(R.id.toolbarMenu))

        // Get the areas
        var sp = com.example.gdsudao.utils.SharedPreferences()
        var areas = sp.RecuperarListaAreas(this)

        val bundle = intent.extras
        val areaIndex = bundle?.getInt("item")

        if(areaIndex != null){
            // Toast.makeText(this, "Area: ${areaIndex}", Toast.LENGTH_SHORT).show(); FUNCIONANDO
            Log.println(Log.DEBUG, "area", "Area: ${areas[areaIndex].toString()}")
        }else{
            Toast.makeText(this, "Erro ao apresentar os detalhes.", Toast.LENGTH_SHORT).show()
        }

    }
}